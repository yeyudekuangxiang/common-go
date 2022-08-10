package cron

import (
	"log"
	ebusiness "mio/internal/pkg/model/entity/business"
	sbusiness "mio/internal/pkg/service/business"
)

func businessCron() {
	AddFunc("0 0 * * ?", func() {
		log.Println("每天执行一次")
		sbusiness.DefaultCarbonRankService.InitUserRank(ebusiness.RankDateTypeDay)

	})

	AddFunc("0 0 * * ?", func() {
		log.Println("每天执行一次")
		sbusiness.DefaultCarbonRankService.InitDepartmentRank(ebusiness.RankDateTypeDay)
	})

	AddFunc("0 0 ? * 1", func() {
		log.Println("每周执行一次")
		sbusiness.DefaultCarbonRankService.InitUserRank(ebusiness.RankDateTypeWeek)
	})

	AddFunc("0 0 ? * 1", func() {
		log.Println("每周执行一次")
		sbusiness.DefaultCarbonRankService.InitDepartmentRank(ebusiness.RankDateTypeWeek)
	})

	AddFunc("0 0 1 1/1 ?", func() {
		log.Println("每月执行一次")
		sbusiness.DefaultCarbonRankService.InitUserRank(ebusiness.RankDateTypeMonth)
	})

	AddFunc("0 0 1 1/1 ?", func() {
		log.Println("每月执行一次")
		sbusiness.DefaultCarbonRankService.InitDepartmentRank(ebusiness.RankDateTypeMonth)
	})
}
