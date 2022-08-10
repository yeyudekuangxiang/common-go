package initialize

import (
	"github.com/shopspring/decimal"
)

func Initialize(configPath string) {
	InitIni(configPath)
	InitConsoleLog()
	InitDB()
	InitRedis()
	InitValidator()
	InitWeapp()
	InitOss()
	InitWxoa()
	decimal.MarshalJSONWithoutQuotes = true
}
