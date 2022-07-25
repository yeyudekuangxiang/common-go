package initialize

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
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
	rand.Seed(time.Now().UnixMilli())
}
