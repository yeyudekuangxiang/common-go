package initialize

import (
	"go.uber.org/zap"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	mzap "mio/pkg/zap"
)

func InitConsoleLog() {
	logger := mzap.DefaultLogger(config.Config.Log.Level).WithOptions(zap.Fields(zap.String("scene", "log"))).Sugar()
	*app.Logger = *logger
}
func InitFileLog() {
	var loggerConfig mzap.LoggerConfig
	var err error
	err = util.MapTo(config.Config.Log, &loggerConfig)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig.Path = "runtime"
	loggerConfig.FileName = "log.log"
	logger := mzap.NewZapLogger(loggerConfig).Sugar()
	*app.Logger = *logger
}
