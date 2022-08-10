package initialize

import (
	"gopkg.in/ini.v1"
	"log"
	"mio/config"
	"mio/internal/pkg/service"
)

func InitIni(source interface{}) {
	f, err := ini.Load(source)
	if err != nil {
		log.Panic(err)
	}
	err = f.MapTo(&config.Config)
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
