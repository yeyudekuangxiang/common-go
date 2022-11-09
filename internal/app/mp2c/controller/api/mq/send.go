package mq

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	messageSrv "mio/internal/pkg/service/message"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/message"
	"mio/pkg/errno"
)

var DefaultMqController = SmsSendController{}

type SmsSendController struct {
}

//发送营销短信

func (c SmsSendController) SendSms(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetSendSmsForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//调用模版服务，获取模版内容
	ctxMio := context.NewMioContext(context.WithContext(ctx.Request.Context()))
	messageService := messageSrv.NewWebMessageService(ctxMio)
	templateContent := messageService.GetTemplate(form.TemplateKey)
	if templateContent == "" {
		return nil, errno.ErrCommon.WithMessage(form.TemplateKey + "该模版不存在")
	}

	//{"code":"0","failNum":"0","successNum":"1","msgId":"22110915322300602201000033772693","time":"20221109153223","errorMsg":""}
	//{"code":"102","msgId":"","time":"20221109153305","errorMsg":"密码错误"}

	body, err := message.SendMarketSms(templateContent, form.Phone, form.Msg)
	if err != nil {
		log.Println("短信发送失败", err, form.Phone, form.Msg)
		return nil, err
	}

	//发送短信
	return gin.H{
		"body":       body,
		"templateId": form.TemplateKey,
		"phone":      form.Phone,
		"msg":        form.Msg,
	}, nil
}

//发送验证码短信

func (c SmsSendController) SendYzmSms(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetSendYzmSmsForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//{"code":"0","failNum":"0","successNum":"1","msgId":"22110915322300602201000033772693","time":"20221109153223","errorMsg":""}
	//{"code":"102","msgId":"","time":"20221109153305","errorMsg":"密码错误"}

	body, err := message.SendYZMSms(form.Phone, form.Code)
	if err != nil {
		log.Println("短信发送失败", err, form.Phone, form.Code)
		return nil, err
	}

	//发送短信
	return gin.H{
		"body":  body,
		"phone": form.Phone,
		"code":  form.Code,
	}, nil
}

//诸葛打点

func (c SmsSendController) SendZhuGe(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetSendZhugeForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	zhuGeAttr := make(map[string]interface{})
	err := json.Unmarshal([]byte(form.Data), &zhuGeAttr)
	if err != nil {
		return nil, err
	}
	//上报到诸葛
	err = track.DefaultZhuGeService().TrackWithErr(form.EventKey, form.Openid, zhuGeAttr)
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
