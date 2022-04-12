package initialize

import (
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	"mio/pkg/db"
)

func InitDB() {
	var conf db.Config
	err := util.MapTo(config.Config.Database, &conf)
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
