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
}
