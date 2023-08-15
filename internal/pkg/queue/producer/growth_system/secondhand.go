package growth_system

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func GrowthSystemMallRedemption(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemMallRedemption] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemMallRedemption},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemMallRedemption] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemSecondHandMarket(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemSecondHandMarket] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemSecondHandMarket},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemSecondHandMarket] method [mq.Publish] failed :%v", err)
	}
	return
}
