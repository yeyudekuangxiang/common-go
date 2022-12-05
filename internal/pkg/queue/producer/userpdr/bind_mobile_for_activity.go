package userpdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/usermsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func BindMobileForActivity(msg usermsg.BindMobile) error {
	return nil
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return app.QueueProduct.Publish(marshal,
		[]string{routerkey.BindMobileForActivity},
		rabbitmq.WithPublishOptionsExchange("userExchange"))
}
