package initialize

import (
	"log"
	"mio/config"
	"mio/core/app"
	"mio/internal/util"
	"mio/internal/zap"
)

func InitLog() {
	var loggerConfig zap.LoggerConfig
	var err error
	err = util.MapTo(config.App.Log, &loggerConfig)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig.Path = "runtime"
	loggerConfig.FileName = "log.log"
	app.Logger = zap.NewZapLogger(loggerConfig).Sugar()
}
