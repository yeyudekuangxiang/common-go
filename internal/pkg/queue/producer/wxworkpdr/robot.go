package wxworkpdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/wxworkqueue"
	"mio/pkg/wxwork"
)

func SendRobotMessage(message wxworkqueue.RobotMessage) error {
	return wxwork.SendRobotMessage(message.Key, message.Message)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return app.QueueProduct.Publish(data, []string{"wxwork.robot"}, rabbitmq.WithPublishOptionsExchange("wxwork"))
}
