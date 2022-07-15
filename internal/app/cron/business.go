package cron

import (
	"mio/internal/pkg/core/app"
	ebusiness "mio/internal/pkg/model/entity/business"
	sbusiness "mio/internal/pkg/service/business"
)

func businessCron() {
	id, err := c.AddFunc("55 14 * * ?", func() {
		app.Logger.Info("每天执行一次")
		sbusiness.DefaultCarbonRankService.InitUserRank(ebusiness.RankDateTypeDay)
	})
	app.Logger.Info(id, err)
	id, err = c.AddFunc("55 14 * * ?", func() {
		app.Logger.Info("每天执行一次")
		sbusiness.DefaultCarbonRankService.InitDepartmentRank(ebusiness.RankDateTypeDay)
	})
	app.Logger.Info(id, err)
	id, err = c.AddFunc("0 0 ? * 2", func() {
		app.Logger.Info("每周执行一次")
		sbusiness.DefaultCarbonRankService.InitUserRank(ebusiness.RankDateTypeWeek)
	})
	app.Logger.Info(id, err)
	id, err = c.AddFunc("0 0 ? * 2", func() {
		app.Logger.Info("每周执行一次")
		sbusiness.DefaultCarbonRankService.InitDepartmentRank(ebusiness.RankDateTypeWeek)
	})
	app.Logger.Info(id, err)
	id, err = c.AddFunc("0 0 1 1/1 ?", func() {
		app.Logger.Info("每月执行一次")
		sbusiness.DefaultCarbonRankService.InitUserRank(ebusiness.RankDateTypeMonth)
	})
	app.Logger.Info(id, err)
	id, err = c.AddFunc("0 0 1 1/1 ?", func() {
		app.Logger.Info("每月执行一次")
		sbusiness.DefaultCarbonRankService.InitDepartmentRank(ebusiness.RankDateTypeMonth)
	})
	app.Logger.Info(id, err)
}
