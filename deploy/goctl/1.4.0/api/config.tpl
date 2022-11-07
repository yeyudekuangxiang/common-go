package config

import {{.authImport}}
import (
    "github.com/zeromicro/go-zero/zrpc"
    "{{.projectPath}}/common/globalclient"
)

type Config struct {
	rest.RestConf
    GlobalClientConf globalclient.GlobalClientConf
	{{.auth}}
	{{.jwtTrans}}
	RpcConf zrpc.RpcClientConf
	JwtAuth       struct {
		AccessSecret string
	}
	Debug bool
}
