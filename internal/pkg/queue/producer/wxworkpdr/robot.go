package wxworkpdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxwork"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/wxworkmsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func SendRobotMessage(message wxworkmsg.RobotMessage) error {
	return wxwork.SendRobotMessage(message.Key, message.Message)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return app.QueueProduct.Publish(data, []string{routerkey.WxWorkRobot}, rabbitmq.WithPublishOptionsExchange("wxwork"))
}
