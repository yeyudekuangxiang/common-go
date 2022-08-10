package main

import (
	"flag"
	"log"
	"mio/internal/app/consumer"
	"mio/internal/pkg/core/initialize"
	"os"
	"os/signal"
)

var (
	flagConf = flag.String("c", "./config.ini", "-c")
)

func init() {
	flag.Parse()
	initialize.Initialize(*flagConf)
}
func main() {
	consumer.Run()
	defer func() {
		log.Println("关闭消费者...")
		err := consumer.Close()
		if err != nil {
			log.Println("关闭消费者异常", err)
		} else {
			log.Println("关闭消费者成功")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
