package initialize

import (
	"github.com/medivhzhan/weapp/v3"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
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
	log.Println("初始化weapp组件...")
	weappSetting := config.Config.Weapp
	c := weapp.NewClient(weappSetting.AppId, weappSetting.Secret, weapp.WithCache(noCache{}))
	*app.Weapp = *c
	log.Println("初始化weapp组件成功")
}
