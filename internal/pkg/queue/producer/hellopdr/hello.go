package hellopdr

import (
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
)

func Hello(msg []byte) error {
	return app.QueueProduct.Publish(msg, []string{"hello"}, rabbitmq.WithPublishOptionsExchange("hello-exchange"))
}
