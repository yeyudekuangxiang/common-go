package main

import (
	"flag"
	"mio/internal/app/cron"
	"mio/internal/app/mp2c/server"
	"mio/internal/pkg/core/initialize"
	"os"
	"os/signal"
)

var (
	flagConf = flag.String("c", "./config.ini", "-c")
)

func init() {
	flag.Parse()

	initialize.Initialize(*flagConf)

	server.InitServer()
}
func main() {
	cron.Run()

	server.RunServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	server.CloseServer()
}
