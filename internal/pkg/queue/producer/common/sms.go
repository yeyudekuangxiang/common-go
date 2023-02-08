package common

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/smsmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

const (
	smsToken = "OSiM9W3dkaSsPDrd1Dllp"
)

func SendSms(message smsmsg.SmsMessage) error {
	marshal, err := json.Marshal(message)
	if err != nil {
		return err
	}
	send := smsmsg.HttpSmsMessage{
		Url:              config.Config.MqArgs.SmsUrl,
		Token:            smsToken,
		Method:           "post",
		ContentType:      "application/json",
		Body:             string(marshal),
		SuccessHttpCodes: []int{200},
	}
	sendBody, err := send.Byte()
	if err != nil {
		return err
	}
	app.Logger.Infof("删除活动发送短信")
	return app.QueueProduct.Publish(sendBody, []string{routerkey.HttpRouterKeys}, rabbitmq.WithPublishOptionsExchange("httpExchange"))
}
