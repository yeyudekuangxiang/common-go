package communitypdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/communitymsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func ReviewImage(msg communitymsg.ReviewImage) error {
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return app.QueueProduct.Publish(marshal,
		[]string{routerkey.ReviewImage},
		rabbitmq.WithPublishOptionsExchange("reviewExchange"))
}
