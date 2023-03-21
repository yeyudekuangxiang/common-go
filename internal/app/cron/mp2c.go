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
	AddFunc("0 1 * * ?", func() {
		log.Println("每天执行一次")
		service.NewCarbonTransactionService(context.NewMioContext()).AddHistory()
	})
	//签到提醒
	AddFunc("0/30 8,9,10,11,12,13,14,15,16,17,18,19,20 * * ?", func() {
		log.Println("8,9,10,11,12,13,14,15,16,17,18,19,20 每30分钟执行一次")
		/*if !util.DefaultLock.Lock("send", time.Minute*5) {
			return
		}*/
		service := messageSrv.MessageService{}
		service.SendMessageToSignUser()
	})
	//打卡挑战提醒
	AddFunc("0 18 * * ?", func() {
		service := messageSrv.MessageService{}
		service.SendMessageToCarbonPk()
	})
}
