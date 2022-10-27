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
	zhuGeAttrV2 := make(map[string]interface{})
	errV6 := json.Unmarshal([]byte(msg.Date), &zhuGeAttrV2)
	if errV6 != nil {
		return rabbitmq.Ack
	}
	//上报到诸葛
	/*zhuGeAttr := make(map[string]interface{}, 0)
	zhuGeAttr["来源"] = "1"
	zhuGeAttr["渠道"] = "3"
	zhuGeAttr["城市code"] = "5"
	zhuGeAttr["openid"] = "4"*/
	track.DefaultZhuGeService().Track(msg.EventKey, msg.Openid, zhuGeAttrV2)
	return rabbitmq.Ack
}
