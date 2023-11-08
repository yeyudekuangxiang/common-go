package growth_system

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func GrowthSystemStep(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemStep] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemTravelStep},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemStep] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemBus(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemBus] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemTravelBus},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemBus] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemRide(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemRide] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemTravelRide},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemRide] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemSubway(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemSubway] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemTravelSubway},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemSubway] method [mq.Publish] failed :%v", err)
	}
	return
}
