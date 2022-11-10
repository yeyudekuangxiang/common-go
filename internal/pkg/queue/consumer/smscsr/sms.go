package smscsr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"mio/internal/pkg/queue/types/message/smsmsg"
	"mio/internal/pkg/util/message"
)

func SendSms(delivery rabbitmq.Delivery) rabbitmq.Action {

	msg := smsmsg.MsgMessage{}
	err := json.Unmarshal(delivery.Body, &msg)
	if err != nil {
		log.Println("转换发送短信失败", err, string(delivery.Body))
		return rabbitmq.Ack
	}
	//发送短信
	message.SendYZM(msg.Phone, msg.Msg)
	return rabbitmq.Ack
}
