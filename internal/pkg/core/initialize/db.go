package initialize

import (
	"go.uber.org/zap"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	"mio/pkg/db"
	mzap "mio/pkg/zap"
)

func InitDB() {
	log.Println("初始化数据库连接...")
	var conf db.Config
	err := util.MapTo(config.Config.Database, &conf)
	if err != nil {
		log.Panic(err)
	}
	conf.Logger = mzap.NewGormLogger(mzap.DefaultLogger(config.Config.Log.Level).WithOptions(zap.Fields(zap.String("scene", "database"))).Sugar())
	gormDb, err := db.NewDB(conf)
	if err != nil {
		log.Panic(err)
	}
	*app.DB = *gormDb
	log.Println("初始化数据库连接成功")
}

func InitBusinessDB() {
	log.Println("初始化企业版数据库连接...")
	var conf db.Config
	err := util.MapTo(config.Config.DatabaseBusiness, &conf)
	if err != nil {
		log.Panic(err)
	}
	conf.Logger = mzap.NewGormLogger(mzap.DefaultLogger(config.Config.Log.Level).WithOptions(zap.Fields(zap.String("scene", "database"))).Sugar())
	gormDb, err := db.NewDB(conf)
	if err != nil {
		log.Panic(err)
	}
	*app.BusinessDB = *gormDb
	log.Println("初始化企业版数据库连接成功")
}
