package main

import (
	"flag"
	"mio/core/initialize"
	"mio/internal/acm"
	"log"
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
		initialize.InitIni(*flagConf)
	case "dev":
		acmClient, err := acm.NewClient(AcmConf)
		if err != nil {
			log.Fatal("create acm client:", err)
		}
		content, err := acmClient.GetConfig(DevAcmGroup, DevAcmDataId)
		if err != nil {
			log.Fatal("get acm config:", err)
		}
		initialize.InitIni([]byte(content))
	case "prod":
		acmClient, err := acm.NewClient(AcmConf)
		if err != nil {
			log.Fatal("create acm client:", err)
		}
		content, err := acmClient.GetConfig(ProdAcmGroup, ProdAcmDataId)
		if err != nil {
			log.Fatal("get acm config:", err)
		}
		initialize.InitIni([]byte(content))
	default:
		log.Fatal("error env:", *flagEnv)
	}
}
func init() {
	flag.Parse()

	initIni()
	initialize.InitLog()
	initialize.InitDB()
	initialize.InitValidator()
	initialize.InitServer()
}
func main() {
	initialize.RunServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	initialize.CloseServer()
}
