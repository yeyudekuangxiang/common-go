package initialize

import (
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	"mio/pkg/zap"
)

func InitLog() {
	var loggerConfig zap.LoggerConfig
	var err error
	err = util.MapTo(config.Config.Log, &loggerConfig)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig.Path = "runtime"
	loggerConfig.FileName = "log.log"
	app.Logger = zap.NewZapLogger(loggerConfig).Sugar()
}
