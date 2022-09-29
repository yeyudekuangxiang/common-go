package cron

import (
	"log"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	messageSrv "mio/internal/pkg/service/message"
)

func mp2cCron() {
	/*	AddFunc("0 0 * * ?", func() {
		log.Println("每天执行一次")
		service.NewCarbonTransactionService(context.NewMioContext()).AddClassify()
	})*/

	//碳历史 每天跑一次
	AddFunc("0 0 * * ?", func() {
		log.Println("每天执行一次")
		service.NewCarbonTransactionService(context.NewMioContext()).AddHistory()
	})

	//签到提醒
	AddFunc("0 0/30 8-9 * ?", func() {
		log.Println("每天执行一次")
		service := messageSrv.MessageService{}
		service.SendMessageToSignUser()
	})

	AddFunc("0 0/30 11-14 * ?", func() {
		log.Println("每天执行一次")
		service := messageSrv.MessageService{}
		service.SendMessageToSignUser()
	})

	AddFunc("0 0/30 17-20 * ?", func() {
		log.Println("每天执行一次")
		service := messageSrv.MessageService{}
		service.SendMessageToSignUser()
	})
}