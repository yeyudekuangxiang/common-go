package initialize

import (
	"github.com/wagslane/go-rabbitmq"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/logger/zap"
)

func initQueueProducer() {
	println(config.Config.AMQP.Url)
	log.Println("初始化amqp生产者...")
	pub, err := rabbitmq.NewPublisher(config.Config.AMQP.Url, rabbitmq.Config{}, rabbitmq.WithPublisherOptionsLogger(zap.NewRabbitmqLogger(app.Logger)))
	if err != nil {
		if config.Config.App.Env == "prod" {
			log.Fatal("初始化amqp生产者失败", err)
		} else {
			log.Println("初始化amqp生产者失败", err)
		}

	} else {
		app.QueueProduct = pub
		log.Println("初始化amqp生产者成功")
	}

	if app.QueueProduct != nil {
		publishConfirm := app.QueueProduct.NotifyPublish()
		returnCh := app.QueueProduct.NotifyReturn()

		returnMsg := make(map[string]int)
		go func() {
			for {
				select {
				case msg := <-returnCh:
					app.Logger.Errorf("消息被退回 %s %+v", msg.MessageId, msg)
					if returnMsg[msg.MessageId] >= 3 {
						app.Logger.Errorf("重试3次失败 %s %+v", msg.MessageId, msg)
						continue
					}
					err = app.QueueProduct.Publish(msg.Body, []string{msg.RoutingKey})
					returnMsg[msg.MessageId] = returnMsg[msg.MessageId] + 1
					if err != nil {
						app.Logger.Errorf("消息重发失败 %s %+v", msg.MessageId, msg)
					}
				}
			}
		}()

		go func() {
			for {
				select {
				case c := <-publishConfirm:
					if !c.Ack {
						app.Logger.Errorf("消息队列发送失败 %d %d", c.ReconnectionCount, c.DeliveryTag)
					} else {
						app.Logger.Info("消息发送成功", c)
					}
				}
			}
		}()
	}
}
func closeQueueProducer() {

	if app.QueueProduct != nil {
		log.Println("关闭QueueProduct")
		err := app.QueueProduct.Close()
		if err != nil {
			log.Println("关闭QueueProduct异常", err)
		} else {
			log.Println("关闭QueueProduct成功")
		}
	}

}
