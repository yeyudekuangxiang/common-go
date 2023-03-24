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
	marshal, err := json.Marshal(topic)
	if err != nil {
		app.Logger.Errorf("城市碳秘 %+v", err.Error())
		return err
	}
	err = app.QueueProduct.Publish(marshal, []string{routerkey.TopicSeekingStore}, rabbitmq.WithPublishOptionsExchange("lvmiao"))
	if err != nil {
		app.Logger.Errorf("城市碳秘 %+v", err.Error())
		return err
	}
	return nil
}
