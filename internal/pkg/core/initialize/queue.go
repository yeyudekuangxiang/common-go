package initialize

import (
	"github.com/wagslane/go-rabbitmq"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/zap"
)

func initQueueProducer() {

	log.Println("初始化amqp生产者...")

	pub, err := rabbitmq.NewPublisher(config.Config.AMQP.Url, rabbitmq.Config{}, rabbitmq.WithPublisherOptionsLogger(zap.NewRabbitmqLogger(app.Logger)))
	if err != nil {

		if config.Config.App.Env == "prod" {
			log.Fatal("初始化amqp生产者失败", err)
		} else {
			log.Println("初始化amqp生产者失败", err)
		}

	} else {
		*app.QueueProduct = *pub
		log.Println("初始化amqp生产者成功")
	}

}
