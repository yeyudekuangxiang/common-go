package wxworkcsr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxwork"
	"log"
	"mio/internal/pkg/queue/types/message/wxworkmsg"
)

func DealWxWorkRobot(delivery rabbitmq.Delivery) rabbitmq.Action {
	robotMsg := wxworkmsg.RobotMessage{}
	err := json.Unmarshal(delivery.Body, &robotMsg)
	if err != nil {
		log.Println("转换企业微信消息机器人失败", err, string(delivery.Body))
		return rabbitmq.Ack
	}

	err = wxwork.SendRobotMessageRaw(robotMsg.Key, robotMsg.Type, robotMsg.Message)
	if err != nil {
		log.Println("发送消息到企业微信机器人失败", err, string(delivery.Body))
	}
	return rabbitmq.Ack
}
