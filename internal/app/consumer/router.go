package consumer

import (
	"mio/internal/pkg/queue/consumer/hellocsr"
	"mio/internal/pkg/queue/consumer/smscsr"
	"mio/internal/pkg/queue/consumer/wxworkcsr"
	"mio/internal/pkg/queue/consumer/zhugecsr"
	"mio/internal/pkg/queue/types/routerkey"
)

func Router() {
	StartConsume("helloqueue", "hello-exchange", []string{routerkey.Hello}, false, hellocsr.DealHello)
	StartConsume("wxworkqueue", "wxwork-exchange", []string{routerkey.WxWorkRobot}, false, wxworkcsr.DealWxWorkRobot)
	StartConsume("smsqueue", "sms-exchange", []string{routerkey.SmsSend}, false, smscsr.SendSms)
	StartConsume("zhugequeue", "zhuge-exchange", []string{routerkey.ZhugeSend}, false, zhugecsr.SendToZhuge)
}
