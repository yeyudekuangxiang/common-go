package consumer

import (
	"mio/internal/pkg/queue/consumer/hellocsr"
	"mio/internal/pkg/queue/consumer/wxworkcsr"
	"mio/internal/pkg/queue/types/routerkey"
)

func Router() {
	StartConsume("helloqueue", "hello-exchange", []string{routerkey.Hello}, false, hellocsr.DealHello)
	StartConsume("wxworkqueue", "wxwork-exchange", []string{routerkey.WxWorkRobot}, false, wxworkcsr.DealWxWorkRobot)
}
