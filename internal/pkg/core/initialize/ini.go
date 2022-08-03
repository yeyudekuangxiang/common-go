package initialize

import (
	"gopkg.in/ini.v1"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
)

func InitIni(source interface{}) {
	f, err := ini.Load(source)
	if err != nil {
		log.Panic(err)
	}
	app.Ini = f
	err = app.Ini.MapTo(&config.Config)
	if err != nil {
		log.Panic(err)
	}
	afterInitIni()
}

func afterInitIni() {
	service.InitDefaultDuibaService()
	service.InitDefaultOssService()
	service.InitDefaultOCRService()
}
