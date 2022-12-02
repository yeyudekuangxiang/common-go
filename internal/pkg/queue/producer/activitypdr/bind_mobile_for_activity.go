package activitypdr

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/routerkey"
)

func BindMobileForActivity(msg []byte) error {
	return app.QueueProduct.Publish(msg,
		[]string{routerkey.ActivityNewUser},
		rabbitmq.WithPublishOptionsExchange("activityExchange"))
}
