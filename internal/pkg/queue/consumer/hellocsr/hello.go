package hellocsr

import (
	"github.com/wagslane/go-rabbitmq"
	"log"
)

func DealHello(delivery rabbitmq.Delivery) rabbitmq.Action {
	log.Println("deal hello1 msg", delivery.RoutingKey, string(delivery.Body))
	return rabbitmq.Ack
}
