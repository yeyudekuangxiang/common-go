package zhugecsr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"log"
	"mio/internal/pkg/queue/types/message/zhugemsg"
	"mio/internal/pkg/service/track"
)

func SendToZhuge(delivery rabbitmq.Delivery) rabbitmq.Action {
	msg := zhugemsg.ZhugeMessage{}
	err := json.Unmarshal(delivery.Body, &msg)
	if err != nil {
		log.Println("转换埋点到诸葛失败", err, string(delivery.Body))
		return rabbitmq.Ack
	}
	zhuGeAttr := make(map[string]interface{})
	err = json.Unmarshal([]byte(msg.Date), &zhuGeAttr)
	if err != nil {
		return rabbitmq.Ack
	}
	//上报到诸葛
	track.DefaultZhuGeService().Track(msg.EventKey, msg.Openid, zhuGeAttr)
	return rabbitmq.Ack
}
