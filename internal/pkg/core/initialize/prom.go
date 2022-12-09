package initialize

import (
	"github.com/zeromicro/go-zero/core/prometheus"
	"log"
	"mio/config"
)

func InitProm() {
	if len(config.Config.Prometheus.Host) != 0 {
		log.Println("初始化prometheus...")
		prometheus.StartAgent(prometheus.Config{
			Host: config.Config.Prometheus.Host,
			Port: config.Config.Prometheus.Port,
			Path: config.Config.Prometheus.Path,
		})
		log.Println("初始化prometheus成功")
	}
}
