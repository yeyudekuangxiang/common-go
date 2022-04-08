package main

import (
	"flag"
	"mio/internal/app/mp2c/router"
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

	router.InitServer()
}
func main() {
	router.RunServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	router.CloseServer()
}
