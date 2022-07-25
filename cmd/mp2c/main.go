package main

import (
	"flag"
	"math/rand"
	"mio/internal/app/cron"
	"mio/internal/app/mp2c/server"
	"mio/internal/pkg/core/initialize"
	"os"
	"os/signal"
	"time"
)

var (
	flagConf = flag.String("c", "./config.ini", "-c")
)

func init() {
	flag.Parse()

	rand.Seed(time.Now().UnixMilli())

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
