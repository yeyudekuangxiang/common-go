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
	log.Println("初始化日志组件...")
	logger := mzap.DefaultLogger(config.Config.Log.Level).WithOptions(zap.Fields(zap.String("scene", "log"))).Sugar()
	*app.Logger = *logger
	log.Println("初始化日志组件成功")
}
func InitFileLog() {
	log.Println("初始化日志组件...")
	var loggerConfig mzap.LoggerConfig
	var err error
	err = util.MapTo(config.Config.Log, &loggerConfig)
	if err != nil {
		log.Panic(err)
	}
	loggerConfig.Path = "runtime"
	loggerConfig.FileName = "log.log"
	logger := mzap.NewZapLogger(loggerConfig).Sugar()
	*app.Logger = *logger
	log.Println("初始化日志组件成功")
}
