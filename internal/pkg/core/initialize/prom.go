package initialize

import (
	"github.com/zeromicro/go-zero/core/prometheus"
	"mio/config"
)

func InitProm() {
	prometheus.StartAgent(prometheus.Config{
		Host: config.Config.Prometheus.Host,
		Port: config.Config.Prometheus.Port,
		Path: config.Config.Prometheus.Path,
	})
}
