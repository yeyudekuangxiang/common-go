package config

import {{.authImport}}
import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	rest.RestConf
	{{.auth}}
	{{.jwtTrans}}
	RpcConf zrpc.RpcClientConf
	JwtAuth       struct {
		AccessSecret string
	}
	Debug bool
}
