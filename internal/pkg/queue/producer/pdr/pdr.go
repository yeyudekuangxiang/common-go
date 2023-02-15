package pdr

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"mio/internal/pkg/core/app"
)

func PublishLogErr(info interface{}, routerKeys []string, exchange string) error {
	data, err := json.Marshal(info)
	if err != nil {
		app.Logger.Errorf("推送数据到rabbitmq异常 %+v %+v %+v %+v\n", string(data), routerKeys, exchange, err)
		return err
	}
	err = app.QueueProduct.Publish(data, routerKeys, rabbitmq.WithPublishOptionsExchange(exchange))
	if err != nil {
		app.Logger.Errorf("推送数据到rabbitmq异常 %+v %+v %+v %+v\n", string(data), routerKeys, exchange, err)
	}
	return err
}
func PublishDataLogErr(data []byte, routerKeys []string, exchange string) error {
	err := app.QueueProduct.Publish(data, routerKeys, rabbitmq.WithPublishOptionsExchange(exchange))
	if err != nil {
		app.Logger.Errorf("推送数据到rabbitmq异常 %+v %+v %+v %+v\n", string(data), routerKeys, exchange, err)
	}
	return err
}
