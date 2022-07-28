package initialize

import (
	"github.com/shopspring/decimal"
)

func Initialize(configPath string) {
	InitIni(configPath)
	InitLog()
	InitDB()
	InitRedis()
	InitValidator()
	InitWeapp()
	initOss()
	decimal.MarshalJSONWithoutQuotes = true
}
