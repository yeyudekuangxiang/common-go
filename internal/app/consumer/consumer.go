package consumer

import (
	"github.com/wagslane/go-rabbitmq"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/zap"
	"os"
)

var consumer rabbitmq.Consumer

func Run() {
	var err error
	consumer, err = rabbitmq.NewConsumer(
		config.Config.AMQP.Url, rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogger(zap.NewRabbitmqLogger(app.Logger)),
	)
	if err != nil {
		log.Fatal("创建消费者失败", err)
	}

	Router()
}
func StartConsume(exchange, queue, topic string, routerKeys []string, handler rabbitmq.Handler) {
	err := consumer.StartConsuming(
		handler,
		queue,
		routerKeys,
		rabbitmq.WithConsumeOptionsBindingExchangeName(exchange),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(topic),
	)
	closeOnErr(err, "创建消费者失败")
}
func closeOnErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)

		log.Println("关闭消费者...")
		err := Close()
		if err != nil {
			log.Println("关闭消费者异常", err)
		} else {
			log.Println("关闭消费者成功")
		}
		os.Exit(1)
	}
}
func Close() error {
	return consumer.Close()
}
