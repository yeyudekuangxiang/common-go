package message

import (
	contextRedis "context"
	"fmt"
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"github.com/pkg/errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"time"
)

type MessageService struct {
}

// SendMiniSubMessage  小程序订阅消息发送
func (srv *MessageService) SendMiniSubMessage(toUser string, page string, template IMiniSubTemplate) (int, error) {
	zhuGeAttr := make(map[string]interface{}, 0) //诸葛打点
	zhuGeAttr["openid"] = toUser
	defer service.DefaultZhuGeService().Track(config.ZhuGeEventName.MessageMiniSubscribe, toUser, zhuGeAttr)

	//发送限制
	redisTemplateKey := fmt.Sprintf(config.RedisKey.MessageLimitByTemplate, time.Now().Format("20060102"))
	redisUserKey := fmt.Sprintf(config.RedisKey.MessageLimitByUser, time.Now().Format("20060102"))
	templateSendCount := app.Redis.ZScore(contextRedis.Background(), redisTemplateKey, toUser).Val()
	userSendCount := app.Redis.ZScore(contextRedis.Background(), redisUserKey, toUser).Val()

	if templateSendCount >= 1 {
		zhuGeAttr["错误码"] = -1
		zhuGeAttr["错误信息"] = "同一模板每人每天最多接收1条消息"
		return -1, errors.New("同一模板每人每天最多接收1条消息")
	}

	if userSendCount >= 2 {
		zhuGeAttr["错误码"] = -2
		zhuGeAttr["错误信息"] = "每人每天最多收到2个不同类型模板消息"
		return -2, errors.New("每人每天最多收到2个不同类型模板消息")
	}

	ret, err := app.Weapp.NewSubscribeMessage().Send(&subscribemessage.SendRequest{
		ToUser:           toUser,
		TemplateID:       template.TemplateId(),
		Page:             page,
		MiniprogramState: subscribemessage.MiniprogramStateDeveloper,
		Data:             template.ToData(),
	})
	if err != nil {
		app.Logger.Info("小程序订阅消息发送失败，http层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, err.Error())
		zhuGeAttr["错误码"] = -3
		zhuGeAttr["错误信息"] = err.Error()
		return -3, err
	}
	zhuGeAttr["错误码"] = ret.ErrCode
	zhuGeAttr["错误信息"] = ret.GetResponseError().Error()
	if ret.ErrCode != 43101 && ret.ErrCode != 0 {
		app.Logger.Info("小程序订阅消息发送失败，业务层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, ret.GetResponseError().Error())
		return ret.ErrCode, ret.GetResponseError()
	}
	app.Redis.ZIncrBy(contextRedis.Background(), redisTemplateKey, 1, toUser).Err() //同一模板每人每天最多接收1条消息
	app.Redis.ZIncrBy(contextRedis.Background(), redisUserKey, 1, toUser).Err()     //每人每天最多收到2个不同类型模板消息
	return ret.ErrCode, nil
}
