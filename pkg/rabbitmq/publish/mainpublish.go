package publish

import (
	"fmt"
	"mio/pkg/rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	kutengone := RabbitMQ.NewRabbitMQRouting("liumeiv2", "liumeiv2")
	for i := 0; i <= 100; i++ {
		kutengone.PublishRouting("Hello kuteng two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
