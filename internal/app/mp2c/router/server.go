package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mio/config"
	"net/http"
	"time"
)

var Server = new(http.Server)

func InitServer() *http.Server {
	//运行模式
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()

	middleware(handler)
	Router(handler)

	*Server = http.Server{
		Handler: handler,
	}
	return Server
}

func RunServer() {
	//gin.DefaultWriter = logger.NewZapLogger(*config.LogConfig)

	configSetting := config.Config.Http
	Server.Addr = fmt.Sprintf(":%d", configSetting.Port)
	Server.ReadTimeout = time.Duration(configSetting.ReadTimeout) * time.Second
	Server.WriteTimeout = time.Duration(configSetting.WriteTimeout) * time.Second
	Server.MaxHeaderBytes = 1 << 20

	//启动
	go func() {
		// 服务连接
		log.Println(fmt.Sprintf("listening: %d", configSetting.Port))
		if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(fmt.Sprintf("listen: %s\n", err))
		}
	}()
}
func CloseServer() {
	var err error
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = Server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
