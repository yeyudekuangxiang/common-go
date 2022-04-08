package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/pugc"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util"
)

func pugcRouter(router *gin.Engine) {
	pugcRouter := router.Group("/pugc")
	pugcRouter.Use(middleware.Throttle())
	//pugcRouter.Use(mustAuth())
	{
		pugcRouter.GET("/addPugc", util.Format(pugc.DefaultPugcController.AddPugc))
		pugcRouter.POST("/ex", util.Format(pugc.DefaultPugcController.ExportExcel))
	}
}
