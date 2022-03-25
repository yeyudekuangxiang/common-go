package initialize

import (
	"log"
	"mio/config"
	"mio/core/app"
	"mio/internal/db"
	"mio/internal/util"
)

func InitDB() {
	var conf db.Config
	err := util.MapTo(config.App.Database, &conf)
	if err != nil {
		log.Fatal(err)
	}
	//创建晓筑规范数据库连接
	gormDb, err := db.NewDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	*app.DB = *gormDb
}
