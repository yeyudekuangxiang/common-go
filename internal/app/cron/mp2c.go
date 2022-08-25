package cron

import (
	"log"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
)

func mp2cCron() {
	AddFunc("0 0 * * ?", func() {
		log.Println("每天执行一次")
		service.NewCarbonTransactionService(context.NewMioContext()).AddClassify()
	})

	AddFunc("0 0 * * ?", func() {
		log.Println("每天执行一次")
		service.NewCarbonTransactionService(context.NewMioContext()).AddHistory()
	})
}
