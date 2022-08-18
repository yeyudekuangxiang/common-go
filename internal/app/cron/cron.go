package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var c = cron.New(cron.WithLocation(time.Local))

func Run() {
	businessCron()
	c.Start()
}
func AddFunc(spec string, f func()) {
	id, err := c.AddFunc(spec, f)
	if err != nil {
		log.Fatal(spec, err)
	}
	log.Println(spec, " next cron time ", c.Entry(id).Schedule.Next(time.Now()))
}
