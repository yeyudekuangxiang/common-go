package point

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"sync"
	"time"
)

type DefaultClientHandle struct {
	ctx          *context.MioContext
	clientHandle ClientHandle
	plugin       clientPlugin
	additional   additional
	paramsMutex  sync.RWMutex // mutex to protect the parameters exposed to the library users
}

type ClientHandle struct {
	UserId       int64
	OpenId       string                      //用户openId
	ImgUrl       string                      //图片地址
	AdminId      int64                       //管理员id
	Type         entity.PointTransactionType //类型
	point        int64                       //变更的积分
	maxPoint     int64
	message      string //记录信息
	bizId        string
	additionInfo string
	identifyImg  map[string]string
}

type ClientOption func(c *ClientHandle)

func ClientOptionNew(option ...ClientOption) *ClientHandle {
	client := &ClientHandle{}
	for _, o := range option {
		o(client)
	}
	return client
}

func WithClientUserId(userId int64) ClientOption {
	return func(handle *ClientHandle) {
		handle.UserId = userId
	}
}

func WithClientOpenId(openId string) ClientOption {
	return func(handle *ClientHandle) {
		handle.OpenId = openId
	}
}

func WithClientImgUrl(imgUrl string) ClientOption {
	return func(handle *ClientHandle) {
		handle.ImgUrl = imgUrl
	}
}

func WithClientType(t entity.PointTransactionType) ClientOption {
	return func(handle *ClientHandle) {
		handle.Type = t
	}
}

func WithClientAdminId(adminId int64) ClientOption {
	return func(handle *ClientHandle) {
		handle.AdminId = adminId
	}
}

//挂件
type clientPlugin struct {
	tracking         *track.ZhuGeService
	pointRepo        *repository.PointRepository
	history          *repository.PointCollectHistoryRepository
	transaction      *repository.PointTransactionRepository
	transactionLimit *repository.PointTransactionCountLimitRepository
}

//额外字段
type additional struct {
	changeType string
	orderId    string //图片识别关键字段
}

func NewClientHandle(ctx *context.MioContext, clientHandle *ClientHandle) *DefaultClientHandle {
	return &DefaultClientHandle{
		ctx:          ctx,
		clientHandle: *clientHandle,
		plugin: clientPlugin{
			tracking:         track.DefaultZhuGeService(),
			pointRepo:        repository.NewPointRepository(ctx),
			history:          repository.NewPointCollectHistoryRepository(ctx),
			transaction:      repository.NewPointTransactionRepository(ctx),
			transactionLimit: repository.NewPointTransactionCountLimitRepository(ctx),
		},
	}
}

// HandleImageCollectCommand 根据图片收集积分
func (c *DefaultClientHandle) HandleImageCollectCommand() (map[string]string, error) {
	cmdDesc := commandMap[string(c.clientHandle.Type)]
	if cmdDesc == nil {
		return nil, errno.ErrRecordNotFound.WithMessage("未找到匹配方法")
	}
	//幂等
	if err := c.checkIdempotency(); err != nil {
		return nil, err
	}
	//检查是否超过次数
	if err := c.checkTimes2(); err != nil {
		//记录日志 返回错误
		return nil, err
	}

	//获取图篇内容
	content, err := c.scanImage(c.clientHandle.ImgUrl)
	if err != nil {
		//记录日志 返回错误
		return nil, err
	}
	imgInfo, err := c.identifyImg(content)
	if err != nil {
		//记录日志 返回错误
		return nil, err
	}
	//添加内容
	marshal, _ := json.Marshal(imgInfo)
	c.withAdditionInfo(string(marshal))
	c.withMessage(fmt.Sprintf("%s", content))
	c.withType(c.clientHandle.Type)
	c.withBizId(util.UUID())
	c.withPoint(cmdDesc.Amount)
	c.withMaxPoint(cmdDesc.MaxAmount)
	//执行function
	if err = c.executeCommandFn(cmdDesc); err != nil {
		//记录日志 返回错误
		return nil, err
	}
	return c.clientHandle.identifyImg, nil
}

// HandlePageDataCommand 收集积分前返回的数据
func (c *DefaultClientHandle) HandlePageDataCommand() (map[string]interface{}, error) {
	pageDataCmd := pageDataMap[string(c.clientHandle.Type)]
	if pageDataCmd == nil {
		return nil, errno.ErrRecordNotFound.WithMessage("未找到匹配方法")
	}
	return pageDataCmd.FnPageData(c)
}

//具体执行方法
func (c *DefaultClientHandle) executeCommandFn(cmdDesc *commandDescription) error {
	defer func() {
		if r := recover(); r != nil {
			//记录错误日志
			//c.writeMessage()
		}
	}()
	if err := cmdDesc.Fn(c); err != nil {
		//记录错误日志
		//c.writeMessage()
		return err
	}
	return nil
}

func (c *DefaultClientHandle) withType(types entity.PointTransactionType) {
	if types != "" {
		c.clientHandle.Type = types
		return
	}
	return
}

func (c *DefaultClientHandle) withPoint(point int64) {
	c.clientHandle.point = point
}

func (c *DefaultClientHandle) withMaxPoint(maxPoint int64) {
	c.clientHandle.maxPoint = maxPoint
}

func (c *DefaultClientHandle) withMessage(message string) {
	if message != "" {
		c.clientHandle.message = message
	}
}

func (c *DefaultClientHandle) withAdminId(adminId int64) {
	if adminId != 0 {
		c.clientHandle.AdminId = adminId
	}
}

func (c *DefaultClientHandle) withBizId(bizId string) {
	if bizId != "" {
		c.clientHandle.bizId = bizId
	}
}

func (c *DefaultClientHandle) withAdditionInfo(additionInfo string) {
	if additionInfo != "" {
		c.clientHandle.additionInfo = additionInfo
	}
}

//保存收集积分记录
func (c *DefaultClientHandle) saveRecord() error {
	history := &entity.PointCollectHistory{
		OpenId: c.clientHandle.OpenId,
		Type:   string(c.clientHandle.Type),
		Info:   c.clientHandle.message,
		Date:   model.Date{Time: time.Now()},
		Time:   model.Time{Time: time.Now()},
	}
	return c.plugin.history.Create(history)
}

// 保存积分
func (c *DefaultClientHandle) savePoint(usrPoint *entity.Point) (int64, error) {
	if err := c.plugin.pointRepo.Save(usrPoint); err != nil {
		return 0, err
	}
	return usrPoint.Balance, nil
}

// 保存积分变动记录 返回积分
func (c *DefaultClientHandle) saveTransAction() (int64, error) {
	pointTransAction := &entity.PointTransaction{
		OpenId:         c.clientHandle.OpenId,
		TransactionId:  c.clientHandle.bizId,
		Type:           c.clientHandle.Type,
		Value:          c.clientHandle.point,
		CreateTime:     model.Time{Time: time.Now()},
		AdditionalInfo: entity.AdditionalInfo(c.clientHandle.additionInfo),
		AdminId:        int(c.clientHandle.AdminId),
		Note:           c.clientHandle.identifyImg["orderId"],
	}
	if err := c.plugin.transaction.Save(pointTransAction); err != nil {
		return 0, err
	}
	return pointTransAction.Value, nil
}

// 保存积分变动次数记录
func (c *DefaultClientHandle) saveTransActionLimit(pointTransActionLimit *entity.PointTransactionCountLimit) error {
	if err := c.plugin.transactionLimit.Save(pointTransActionLimit); err != nil {
		return err
	}
	return nil
}

//更新积分变动次数记录
func (c *DefaultClientHandle) updateTransActionLimit(pointTransActionLimit entity.PointTransactionCountLimit) error {
	if err := c.plugin.transactionLimit.Save(&pointTransActionLimit); err != nil {
		return err
	}
	return nil
}

// 获取用户积分信息
func (c *DefaultClientHandle) findByOpenId() (*entity.Point, error) {
	if c.clientHandle.OpenId == "" {
		return nil, errno.ErrUserNotFound.WithErrMessage("用户未授权")
	}
	p := c.plugin.pointRepo.FindBy(repository.FindPointBy{OpenId: c.clientHandle.OpenId})
	return &p, nil
}

// 增加积分，返回现有积分
func (c *DefaultClientHandle) incPoint(num int64) (int64, error) {
	pointSrv := service.NewPointService(c.ctx)
	balance, err := pointSrv.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       c.clientHandle.OpenId,
		Type:         c.clientHandle.Type,
		BizId:        c.clientHandle.bizId,
		ChangePoint:  c.clientHandle.point,
		AdminId:      int(c.clientHandle.AdminId),
		Note:         c.clientHandle.identifyImg["orderId"],
		AdditionInfo: c.clientHandle.additionInfo,
	})
	return balance, err
}

// 消耗积分，返回现有积分
func (c *DefaultClientHandle) decPoint(num int64) (int64, error) {
	if num == 0 {
		return 0, nil
	}
	if num > 0 {
		num = -num
	}
	c.additional.changeType = "dec"
	usrPoint, err := c.plugin.pointRepo.FindForUpdate(c.clientHandle.OpenId)
	if err != nil {
		return 0, err
	}
	err = c.checkUsrPoints(usrPoint.Balance)
	if err != nil {
		return 0, err
	}
	usrPoint.Balance -= c.clientHandle.point
	point, err := c.savePoint(&usrPoint)
	if err != nil {
		return 0, err
	}
	//添加记录
	err = c.saveRecord()
	if err != nil {
		return 0, err
	}
	return point, nil
}

func (c *DefaultClientHandle) changePointByAdmin() error {
	return nil
}

//更新次数
func (c *DefaultClientHandle) changeLimit() error {
	//获取次数记录
	result := c.plugin.transactionLimit.FindBy(repository.FindPointTransactionCountLimitBy{
		OpenId:          c.clientHandle.OpenId,
		TransactionType: c.clientHandle.Type,
		TransactionDate: model.Date{Time: time.Now()},
	})
	result.CurrentCount += 1
	err := c.saveTransActionLimit(&result)
	if err != nil {
		return err
	}
	return nil
}
