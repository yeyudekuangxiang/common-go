package consumer

import (
	"github.com/wagslane/go-rabbitmq"
	"log"
	"mio/config"
	"mio/internal/app/consumer/hello"
	"mio/internal/pkg/core/app"
	"mio/pkg/zap"
)

func fatalOnErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

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

	err = consumer.StartConsuming(
		hello.DealHello,
		"hello",
		[]string{"hello"},
		rabbitmq.WithConsumeOptionsBindingExchangeName("hello-exchange"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
	)
	fatalOnErr(err, "开始消费失败")

}

func Close() error {
	return consumer.Close()
}
