package mq

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/message"
	"mio/pkg/errno"
)

var DefaultMqController = SmsSendController{}

type SmsSendController struct {
}

func (c SmsSendController) SendSms(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetSendSmsForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	if form.Phone == "" || form.Msg == "" {
		return nil, errno.ErrCommon.WithMessage("手机号或者msg为空")
	}
	body, err := message.SendYZMV5(form.TemplateId, form.Phone, form.Msg)

	err = message.SendSmsV2(form.Phone, form.Msg)
	if err != nil {
		log.Println("短信发送失败", err, form.Phone, form.Msg)
		return nil, err
	}

	//发送短信
	return gin.H{
		"body":       body,
		"templateId": form.TemplateId,
		"phone":      form.Phone,
		"msg":        form.Msg,
	}, nil
}

func (c SmsSendController) SendZhuGe(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetSendZhugeForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	zhuGeAttr := make(map[string]interface{})
	err := json.Unmarshal([]byte(form.Data), &zhuGeAttr)

	/*	zhuGeAttr["phone"] = "18840853003"
		zhuGeAttr["err"] = "err"
		a, _ := json.Marshal(zhuGeAttr)
		println(a)*/

	if err != nil {
		return gin.H{}, nil
	}
	//上报到诸葛
	track.DefaultZhuGeService().Track(form.EventKey, form.Openid, zhuGeAttr)
	return gin.H{}, nil
}
