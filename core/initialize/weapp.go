package initialize

import (
	"github.com/medivhzhan/weapp/v3"
	"log"
	"mio/core/app"
	"time"
)

type weappConfig struct {
	AppId  string
	Secret string
}

type noCache struct {
}

func (noCache) Set(key string, val interface{}, timeout time.Duration) {

}

func (noCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func InitWeapp() {
	conf := weappConfig{}
	if err := app.Ini.Section("weapp").MapTo(&conf); err != nil {
		log.Panic("获取小程序配置文件失败", err)
	}
	app.Weapp = weapp.NewClient(conf.AppId, conf.Secret, weapp.WithCache(noCache{}))
}
