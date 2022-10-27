package wxworkpdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/wxworkqueue"
	"mio/internal/pkg/queue/types/routerkey"
)

func SendRobotMessage(message wxworkqueue.RobotMessage) error {
	//return wxwork.SendRobotMessage(message.Key, message.Message)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return app.QueueProduct.Publish(data, []string{routerkey.WxWorkRobot}, rabbitmq.WithPublishOptionsExchange("wxwork"))
}
