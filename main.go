package main

import (
	"flag"
	"log"
	"mio/core/initialize"
	"mio/internal/acm"
	"os"
	"os/signal"
)

var (
	//env(local,dev,prod) 等于local时是用本地配置文件、dev时是用acm测试配置、prod时是用acm正式配置
	flagEnv  = flag.String("env", "local", "-env")
	flagConf = flag.String("conf", "./config.ini", "-c")
)
var AcmConf = acm.Config{
	Endpoint:    "acm.aliyun.com",
	NamespaceId: "",
	AccessKey:   "",
	SecretKey:   "",
	LogDir:      "acm",
}

const (
	DevAcmGroup   = "DEFAULT_GROUP"
	DevAcmDataId  = ""
	ProdAcmGroup  = "DEFAULT_GROUP"
	ProdAcmDataId = ""
)

func initIni() {
	switch *flagEnv {
	case "local":
		initialize.InitIni("./config-dev.ini")
	case "dev":
		initialize.InitIni("./config-dev.ini")
	case "prod":
		initialize.InitIni("./config-prod.ini")
	default:
		log.Fatal("error env:", *flagEnv)
	}
}
func init() {
	flag.Parse()

	initIni()
	//initialize.InitLog()
	initialize.InitDB()
	initialize.InitValidator()
	initialize.InitWeapp()
	initialize.InitServer()
}
func main() {
	initialize.RunServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	initialize.CloseServer()
}
