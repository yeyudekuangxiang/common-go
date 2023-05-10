package initialize

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/core/app"
)

func Initialize(configPath string) {
	//InitIni(configPath)
	InitYaml(configPath)
	InitLog()
	InitHttpToolLog()
	InitDB()
	InitBusinessDB()
	InitActivityDB()
	InitRedis()
	InitValidator()
	InitOss()
	InitSts()
	InitWxoa()
	InitProm()
	InitRpc()
	InitWeapp(app.RpcService.TokenCenterRpcSrv)
	initQueueProducer()
	decimal.MarshalJSONWithoutQuotes = true
	InitSensors()
}
