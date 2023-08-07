package growth_system

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func GrowthSystemRERecharge(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemRERecharge] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemRERecharge},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemRERecharge] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemRERechargeMio(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemRERechargeMio] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemRERechargeMio},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemRERechargeMio] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemREBatterySwapping(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemREBatterySwapping] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemREBatterySwapping},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemREBatterySwapping] method [mq.Publish] failed :%v", err)
	}
	return
}
