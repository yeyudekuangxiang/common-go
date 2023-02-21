package consumer

import (
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/common-go/logger/zap"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
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

// StartConsume  路由模式默认topic
// queue 消费队列名称
// exchange 交换机名称
// routerKeys 路由key
// durable 消息是否持久化 重要的消息开启持久化
// handler 消费回调
func StartConsume(queue, exchange string, routerKeys []string, durable bool, handler rabbitmq.Handler) {
	options := make([]func(*rabbitmq.ConsumeOptions), 0)
	options = append(options, rabbitmq.WithConsumeOptionsBindingExchangeName(exchange),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"))
	if durable {
		options = append(options, rabbitmq.WithConsumeOptionsBindingExchangeDurable,
			rabbitmq.WithConsumeOptionsQueueDurable)
	}

	err := consumer.StartConsuming(
		handler,
		queue,
		routerKeys,
		options...,
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
