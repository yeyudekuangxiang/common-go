package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/config"
	"mio/core/app"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.App.Redis.Host, config.App.Redis.Port),
		Password: config.App.Redis.Password, // no password set
		DB:       config.App.Redis.DB,       // use default DB
	})
	*app.Redis = *rdb
}
