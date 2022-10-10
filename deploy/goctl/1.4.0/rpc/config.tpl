package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Database DbConf
	Cache    cache.CacheConf
	Debug bool
}
type RedisConf struct {
	Host     string
	Password string `json:",optional"`
	DB       int    `json:",optional"`
}
type DbConf struct {
	Type         string
	Host         string
	UserName     string
	Password     string
	Database     string
	Port         int
	TablePrefix  string `json:",optional"`
	MaxOpenConns int    `json:",optional"` //最大连接数 <=0表示不限制连接数
	MaxIdleConns int    `json:",optional"` //最大空闲数 <=0表示不保留空闲连接
	MaxLifetime  int    `json:",optional"` //连接可重用时间 <=0表示永远可用(单位秒)
}
