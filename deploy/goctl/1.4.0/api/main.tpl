package main

import (
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v4/request"
    "{{.projectPath}}/common/errno"
	"{{.projectPath}}/common/globalclient"

	{{.importPackages}}
	"{{.projectPath}}/common/result"
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func init() {
	request.AuthorizationHeaderExtractor.Extractor = request.HeaderExtractor{"Authorization", "token"}
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	fmt.Printf("配置文件 %+v\n",c)
	errno.Debug = c.Debug

	server := rest.MustNewServer(c.RestConf, rest.WithCors(), rest.WithUnauthorizedCallback(result.HttpUnAuthCallback))
	defer server.Stop()

	globalclient.InitGlobalClient(c)
	defer globalclient.Close()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
