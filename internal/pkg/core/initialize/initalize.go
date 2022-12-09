package initialize

import (
	"github.com/shopspring/decimal"
)

func Initialize(configPath string) {
	InitIni(configPath)
	InitLog()
	InitDB()
	InitBusinessDB()
	InitRedis()
	InitValidator()
	InitWeapp()
	InitOss()
	InitSts()
	InitWxoa()
	InitProm()
	//InitRpc()
	//initQueueProducer()
	decimal.MarshalJSONWithoutQuotes = true
}
