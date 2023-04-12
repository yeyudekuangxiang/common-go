package community

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/communitymsg"
	"mio/internal/pkg/queue/types/routerkey"
)

//SeekingStore 探店
func SeekingStore(topic communitymsg.Topic) error {
	//查询是否为探店标签
	var flag bool
	for _, tag := range topic.Tags {
		if tag.Name == "城市碳秘" {
			flag = true
		}
	}
	if !flag {
		return nil
	}
	marshal, err := json.Marshal(topic)
	if err != nil {
		app.Logger.Errorf("[城市碳秘] json_encode 错误: %+v\n", err.Error())
		return err
	}
	err = app.QueueProduct.Publish(marshal, []string{routerkey.TopicSeekingStore}, rabbitmq.WithPublishOptionsExchange("lvmio"))
	if err != nil {
		app.Logger.Errorf("[城市碳秘] 错误:%+v\n", err.Error())
		return err
	}
	return nil
}
