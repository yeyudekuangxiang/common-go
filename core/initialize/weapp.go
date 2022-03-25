package initialize

import (
	"github.com/medivhzhan/weapp/v3"
	"mio/config"
	"mio/core/app"
	"time"
)

type noCache struct {
}

func (noCache) Set(key string, val interface{}, timeout time.Duration) {

}

func (noCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func InitWeapp() {
	weappSetting := config.App.Weapp
	c := weapp.NewClient(weappSetting.AppId, weappSetting.Secret, weapp.WithCache(noCache{}))
	*app.Weapp = *c
}
