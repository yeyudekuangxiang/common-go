package recieve

import "mio/pkg/rabbitmq/RabbitMQ"

func main() {
	kutengone := RabbitMQ.NewRabbitMQRouting("kuteng", "kuteng_one")
	kutengone.RecieveRouting()
}
