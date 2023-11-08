package carbon

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	carbonmsg "mio/internal/pkg/queue/types/message/carbon"
	"mio/internal/pkg/queue/types/routerkey"
)

func ChangeSuccessToQueue(params carbonmsg.CarbonChangeSuccess) error {
	marshal, err := json.Marshal(params)
	if err != nil {
		return err
	}
	err = app.QueueProduct.Publish(
		marshal,
		[]string{routerkey.CarbonChangeSuccess},
		rabbitmq.WithPublishOptionsExchange("lvmio"),
	)
	if err != nil {
		return err
	}
	return nil
}
