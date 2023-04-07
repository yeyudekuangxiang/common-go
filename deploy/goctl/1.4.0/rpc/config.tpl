package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
    "{{.projectPath}}/common/globalclient"
    "{{.projectPath}}/common/tool/db"
)

type Config struct {
	zrpc.RpcServerConf
	GlobalClientConf globalclient.GlobalClientConf
	Database db.DatabaseConf
	Cache    cache.CacheConf
	Debug bool
}

