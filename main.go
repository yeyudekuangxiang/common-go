package main

import (
	"flag"
	"mio/core/initialize"
	"os"
	"os/signal"
)

var (
	//env(local,dev,prod) 等于local时是用本地配置文件、dev时是用acm测试配置、prod时是用acm正式配置
	flagEnv  = flag.String("env", "local", "-env")
	flagConf = flag.String("c", "./config.ini", "-c")
)

func initIni() {
	initialize.InitIni(*flagConf)
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
