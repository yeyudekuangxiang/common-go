package consumer

import (
	"mio/internal/pkg/queue/consumer/hellocsr"
	"mio/internal/pkg/queue/consumer/wxworkcsr"
)

func Router() {
	StartConsume("hello-exchange", "hello", "topic", []string{"hello"}, hellocsr.DealHello)
	StartConsume("wxwork", "wxworkrobot", "topic", []string{"wxwork.robot"}, wxworkcsr.DealWxWorkRobot)
}
