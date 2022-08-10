package point

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"strings"
)

type clientHandle struct {
	OpenId       string      //用户openId
	ImgUrl       string      //图片地址
	AdminId      int64       //管理员id
	Type         CollectType //类型
	Point        int64       //变更的积分
	Message      string      //记录信息
	BizId        string
	AdditionInfo string
}

type Options struct {
	f func(*clientHandle)
}

func NewClientHandle(openId, imgUrl string) *clientHandle {
	return &clientHandle{
		OpenId: openId,
		ImgUrl: imgUrl,
	}
}

func (c *clientHandle) HandleCommand(types string) {
	types = strings.ToUpper(types)
	cmdDesc := commandMap[types]
	if cmdDesc == nil {
		//记录log 返回错误
	}
	if !cmdDesc.Open {
		//方法未实现 记录log 返回错误
	}
	c.WithType(types)
	c.executeCommandFn(cmdDesc)
}

func (c *clientHandle) executeCommandFn(cmdDesc *CommandDescription) {
	defer func() {
		if r := recover(); r != nil {
			//记录错误日志
		}
		//记录日志
	}()
	if err := cmdDesc.Fn(c); err != nil {
		//记录错误日志
	}
}

func (c *clientHandle) WithType(types string) {
	if types != "" {
		c.Type = CollectType(types)
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
		c.BizId = util.UUID()
	}
}

func (c *clientHandle) WithAdditionInfo(additionInfo string) {
	if additionInfo != "" {
		c.AdditionInfo = additionInfo
	}
}

func (c *clientHandle) saveRecord() error {
	history := &entity.PointCollectHistory{
		OpenId: c.OpenId,
		Type:   string(c.Type),
		Info:   c.Message,
		Date:   model.Date{},
		Time:   model.Time{},
	}
	return repository.DefaultPointCollectHistoryRepository.Create(history)
}

func (c *clientHandle) findByUserId(userId int64) (*entity.Point, error) {
	return nil, nil
}

func (c *clientHandle) findByOpenId() (*entity.Point, error) {
	return nil, nil
}

func (c *clientHandle) incPoint() error {
	return nil
}

func (c *clientHandle) decPoint() error {
	return nil
}

func (c *clientHandle) changePoint() error {
	return nil
}

func (c *clientHandle) changePointByAdmin() error {
	return nil
}

func (c *clientHandle) trackPoint(err error) {
	if err != nil {
		//do something...
	}
	//do something...
}

func (c *clientHandle) writeMessage(code int, message string) {

}
