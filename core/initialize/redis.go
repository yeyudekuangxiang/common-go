package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/config"
	"mio/core/app"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password, // no password set
		DB:       config.Config.Redis.DB,       // use default DB
	})
	*app.Redis = *rdb
}
