package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/core/app"
)

type redisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func InitRedis() {
	c := redisConfig{}
	if err := app.Ini.Section("redis").MapTo(&c); err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})
	*app.Redis = *rdb
}
