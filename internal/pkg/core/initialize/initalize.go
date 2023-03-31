package initialize

import (
	"github.com/shopspring/decimal"
)

func Initialize(configPath string) {
	InitIni(configPath)
	InitLog()
	InitHttpToolLog()
	InitDB()
	InitBusinessDB()
	InitActivityDB()
	InitRedis()
	InitValidator()
	InitWeapp()
	InitOss()
	InitSts()
	InitWxoa()
	InitProm()
	InitRpc()
	initQueueProducer()
	decimal.MarshalJSONWithoutQuotes = true
}
