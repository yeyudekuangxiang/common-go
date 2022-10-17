package initialize

import (
	"github.com/shopspring/decimal"
)

func Initialize(configPath string) {
	InitIni(configPath)
	InitConsoleLog()
	InitDB()
	InitBusinessDB()
	InitRedis()
	InitValidator()
	InitWeapp()
	InitOss()
	InitWxoa()
	InitRpc()
	initQueueProducer()
	decimal.MarshalJSONWithoutQuotes = true
}
