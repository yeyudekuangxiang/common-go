package wxworkpdr

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/wxwork"
	"mio/internal/pkg/queue/types/message/wxworkmsg"
)

func SendRobotMessage(message wxworkmsg.RobotMessage) error {
	return wxwork.SendRobotMessage(message.Key, message.Message)
	/*data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return app.QueueProduct.Publish(data, []string{routerkey.WxWorkRobot}, rabbitmq.WithPublishOptionsExchange("wxwork"))*/
}
