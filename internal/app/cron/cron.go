package cron

import (
	"github.com/robfig/cron/v3"
	"time"
)

var c = cron.New(cron.WithLocation(time.Local))

func Run() {
	businessCron()
	c.Start()
}
