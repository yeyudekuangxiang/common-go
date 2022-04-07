package initialize

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mio/config"
	"mio/core/app"
	"mio/server"
	"net/http"
	"time"
)

func InitServer() *http.Server {
	//运行模式
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()
	*app.Server = http.Server{
		Handler: handler,
	}
	server.Middleware(handler)
	server.Router(handler)
	return app.Server
}

func RunServer() {
	//gin.DefaultWriter = logger.NewZapLogger(*config.LogConfig)

	configSetting := config.Config.Http
	app.Server.Addr = fmt.Sprintf(":%d", configSetting.Port)
	app.Server.ReadTimeout = time.Duration(configSetting.ReadTimeout) * time.Second
	app.Server.WriteTimeout = time.Duration(configSetting.WriteTimeout) * time.Second
	app.Server.MaxHeaderBytes = 1 << 20

	//启动
	go func() {
		// 服务连接
		log.Println(fmt.Sprintf("listening: %d", configSetting.Port))
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(fmt.Sprintf("listen: %s\n", err))
		}
	}()
}
func CloseServer() {
	var err error
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = app.Server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
