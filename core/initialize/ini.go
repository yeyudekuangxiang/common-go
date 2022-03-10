package initialize

import (
	"gopkg.in/ini.v1"
	"log"
	"mio/config"
	"mio/core/app"
)

func InitIni(source interface{}) {
	f, err := ini.Load(source)
	if err != nil {
		log.Fatal(err)
	}
	app.Ini = f
	err = app.Ini.MapTo(&config.App)
	if err != nil {
		log.Fatal(err)
	}
}
