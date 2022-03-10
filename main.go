package main

import (
	"flag"
	"mio/core/initialize"
	"os"
	"os/signal"
)

var (
	flagConf = flag.String("c", "./config.ini", "-c")
)

func initIni() {
	initialize.InitIni(*flagConf)
}
func init() {
	flag.Parse()

	initIni()
	initialize.InitLog()
	initialize.InitDB()
	initialize.InitRedis()
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
