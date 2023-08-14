package growth_system

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func GrowthSystemPublicWelfare(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemPublicWelfare] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemPublicWelfare},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemPublicWelfare] method [mq.Publish] failed :%v", err)
	}
	return
}
