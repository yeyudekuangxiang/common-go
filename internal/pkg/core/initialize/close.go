package initialize

import (
	"log"
	"mio/internal/pkg/core/app"
)

func Close() {
	log.Println("关闭QueueProduct")
	err := app.QueueProduct.Close()
	if err != nil {
		log.Println("关闭QueueProduct异常", err)
	}
	log.Println("关闭QueueProduct成功")
}
