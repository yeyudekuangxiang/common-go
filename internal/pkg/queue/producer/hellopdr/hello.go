package hellopdr

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/routerkey"
)

func Hello(msg []byte) error {
	return app.QueueProduct.Publish(msg, []string{routerkey.Hello}, rabbitmq.WithPublishOptionsExchange("hello-exchange"))
}
