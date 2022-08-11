package consumer

import (
	"mio/internal/pkg/queue/consumer/hellocsr"
	"mio/internal/pkg/queue/consumer/wxworkcsr"
)

func Router() {
	StartConsume("hello-exchange", "hello", []string{"hello"}, false, hellocsr.DealHello)
	StartConsume("wxwork", "wxworkrobot", []string{"wxwork.robot"}, false, wxworkcsr.DealWxWorkRobot)
}
