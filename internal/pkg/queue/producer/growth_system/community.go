package growth_system

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func GrowthSystemCommunityLike(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityLike] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemCommunityLike},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityLike] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemCommunityShare(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityShare] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemCommunityShare},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityShare] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemCommunityComment(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityComment] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemCommunityComment},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityComment] method [mq.Publish] failed :%v", err)
	}
	return
}

func GrowthSystemCommunityPush(param growthsystemmsg.GrowthSystemParam) {
	params, err := checkParams(param)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityPush] method [checkParams] failed :%v", err)
		return
	}

	err = app.QueueProduct.Publish(
		params,
		[]string{routerkey.GrowthSystemCommunityPush},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		app.Logger.Errorf("server [mp2c-go] middleware [GrowthSystemCommunityPush] method [mq.Publish] failed :%v", err)
	}
	return
}
