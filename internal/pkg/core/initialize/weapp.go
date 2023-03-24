package initialize

import (
	"github.com/medivhzhan/weapp/v3/logger"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util/factory"
)

//debug info warn error
var logLevelMap = map[string]logger.Level{
	"debug": logger.Info,
	"info":  logger.Info,
	"warn":  logger.Warn,
	"error": logger.Error,
}

func InitWeapp(tokenCenterRpc tokencenterclient.TokenCenter) {
	log.Println("初始化weapp组件...")
	client, err := factory.NewWxAppFromTokenCenterRpc("lvmioweapp", tokenCenterRpc, logger.Info)
	if err != nil {
		log.Panic(err)
	}
	*app.Weapp = *client
	log.Println("初始化weapp组件成功")
}
