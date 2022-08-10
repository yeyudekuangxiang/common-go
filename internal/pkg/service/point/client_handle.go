package point

import (
	"fmt"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strings"
	"sync"
)

type clientHandle struct {
	ctx          *context.MioContext
	OpenId       string      //用户openId
	ImgUrl       string      //图片地址
	AdminId      int64       //管理员id
	Type         CollectType //类型
	Point        int64       //变更的积分
	Message      string      //记录信息
	BizId        string
	AdditionInfo string
	plugin       clientPlugin
	additional   additional
	paramsMutex  sync.RWMutex // mutex to protect the parameters exposed to the library users
}

//挂件
type clientPlugin struct {
	tracking  *service.ZhuGeService
	pointRepo repository.PointRepository
}

//额外字段
type additional struct {
	changeType string
	orderId    string //图片识别关键字段
}

type Options struct {
	f func(*clientHandle)
}

func NewClientHandle(ctx *context.MioContext, openId, imgUrl string) *clientHandle {
	return &clientHandle{
		ctx:    ctx,
		OpenId: openId,
		ImgUrl: imgUrl,
		plugin: clientPlugin{
			tracking:  service.DefaultZhuGeService(),
			pointRepo: repository.NewPointRepository(ctx),
		},
	}
}

func (c *clientHandle) HandleCollectCommand(types string) error {
	types = strings.ToUpper(types)
	cmdDesc := commandMap[types]
	if cmdDesc == nil {
		//记录log 返回错误
	}
	//检查是否超过次数
	if err := c.checkTimes(cmdDesc.Times); err != nil {
		//记录日志 返回错误
		return err
	}
	//获取图片内容
	intersect, err := c.scanImage(c.ImgUrl)
	if err != nil {
		//记录日志 返回错误
		return err
	}
	c.identifyImg(intersect)
	//幂等
	if err = c.checkIdempotency(); err != nil {
		return err
	}
	//检测当日次数
	if err = c.checkTimes(cmdDesc.Times); err != nil {
		return err
	}
	//添加内容
	c.WithAdditionInfo(fmt.Sprintf("%s", intersect))
	c.WithType(CollectType(types))
	c.WithBizId(util.UUID())
	c.WithPoint(cmdDesc.Amount)
	//执行function
	if err = c.executeCommandFn(cmdDesc); err != nil {
		//记录日志 返回错误
		return err
	}
	return nil
}

func (c *clientHandle) executeCommandFn(cmdDesc *CommandDescription) error {
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

func (c *clientHandle) WithType(types CollectType) {
	if types != "" {
		c.Type = types
	}
}

func (c *clientHandle) WithPoint(point int64) {
	if point != 0 {
		c.Point = point
	} else if c.Type != "" {

	}
}
func (c *clientHandle) WithMessage(message string) {
	if message != "" {
		c.Message = message
	}
}

func (c *clientHandle) WithAdminId(adminId int64) {
	if adminId != 0 {
		c.AdminId = adminId
	}
}

func (c *clientHandle) WithBizId(bizId string) {
	if bizId != "" {
		c.BizId = bizId
	}
}

func (c *clientHandle) WithAdditionInfo(additionInfo string) {
	if additionInfo != "" {
		c.AdditionInfo = additionInfo
	}
}

//保存收集积分记录
func (c *clientHandle) saveRecord() error {
	history := &entity.PointCollectLog{
		OpenId: c.OpenId,
		Type:   string(c.Type),
		Info:   c.Message,
		Point:  c.Point,
		Date:   model.Date{},
		Time:   model.Time{},
	}
	if c.additional.orderId != "" {
		history.AdditionalOrder = c.additional.orderId
	}
	return repository.DefaultPointCollectHistoryRepository.CreateLog(history)
}

// 保存积分
func (c *clientHandle) savePoint(usrPoint *entity.Point) (int64, error) {
	if err := c.plugin.pointRepo.Save(usrPoint); err != nil {
		return 0, err
	}
	return usrPoint.Balance, nil
}

// 获取用户积分信息
func (c *clientHandle) findByOpenId() (*entity.Point, error) {
	if c.OpenId == "" {
		return nil, errno.ErrUserNotFound.WithErrMessage("用户未授权")
	}
	p := c.plugin.pointRepo.FindBy(repository.FindPointBy{OpenId: c.OpenId})
	return &p, nil
}

// 增加积分，返回现有积分
func (c *clientHandle) incPoint(num int64) (int64, error) {
	if num == 0 {
		return 0, nil
	}
	c.additional.changeType = "inc"
	usrPoint, err := c.plugin.pointRepo.FindForUpdate(c.OpenId)
	if err != nil {
		return 0, err
	}
	if usrPoint.Id == 0 {
		usrPoint.OpenId = c.OpenId
		usrPoint.Balance = c.Point
	} else {
		usrPoint.Balance += c.Point
	}
	//操作积分
	point, err := c.savePoint(&usrPoint)
	if err != nil {
		return 0, err
	}
	return point, nil
}

// 消耗积分，返回现有积分
func (c *clientHandle) decPoint(num int64) (int64, error) {
	if num == 0 {
		return 0, nil
	}
	if num > 0 {
		num = -num
	}
	c.additional.changeType = "dec"
	usrPoint, err := c.plugin.pointRepo.FindForUpdate(c.OpenId)
	if err != nil {
		return 0, err
	}
	err = c.checkUsrPoints(usrPoint.Balance)
	if err != nil {
		return 0, err
	}
	usrPoint.Balance -= c.Point
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

func (c *clientHandle) changePointByAdmin() error {
	return nil
}
