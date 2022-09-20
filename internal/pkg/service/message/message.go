package message

import (
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
)

type MessageService struct {
}

func (srv *MessageService) SendMiniSubMessage(toUser string, page string, template IMiniSubTemplate) (int, error) {
	zhuGeAttr := make(map[string]interface{}, 0) //诸葛打点

	ret, err := app.Weapp.NewSubscribeMessage().Send(&subscribemessage.SendRequest{
		ToUser:           toUser,
		TemplateID:       template.TemplateId(),
		Page:             page,
		MiniprogramState: subscribemessage.MiniprogramStateDeveloper,
		Data:             template.ToData(),
	})
	if err != nil {
		app.Logger.Info("小程序订阅消息发送失败，http层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, err.Error())
		zhuGeAttr["openid"] = toUser
		zhuGeAttr["错误码"] = 1
		zhuGeAttr["错误信息"] = err.Error()
		return 1, err
	}
	zhuGeAttr["openid"] = toUser
	zhuGeAttr["错误码"] = ret.ErrCode
	zhuGeAttr["错误信息"] = ret.GetResponseError().Error()
	if ret.ErrCode != 43101 && ret.ErrCode != 0 {
		app.Logger.Info("小程序订阅消息发送失败，业务层，模版%s，toUser%s，错误信息%s", template.TemplateId(), toUser, ret.GetResponseError().Error())
		return ret.ErrCode, ret.GetResponseError()
	}
	service.DefaultZhuGeService().Track(config.ZhuGeEventName.MessageMiniSubscribe, toUser, zhuGeAttr)
	return ret.ErrCode, nil
}
