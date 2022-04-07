package initialize

import (
	"gopkg.in/ini.v1"
	"log"
	"mio/config"
	"mio/core/app"
	"mio/service"
)

func InitIni(source interface{}) {
	f, err := ini.Load(source)
	if err != nil {
		log.Fatal(err)
	}
	app.Ini = f
	err = app.Ini.MapTo(&config.Config)
	if err != nil {
		log.Fatal(err)
	}
	afterInitIni()
}

func afterInitIni() {
	service.InitDefaultDuibaService()
}
