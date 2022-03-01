package server

import (
	"github.com/gin-gonic/gin"
	"mio/controller/pugc"
)

func pugcRouter(router *gin.Engine) {
	pugcRouter := router.Group("/pugc")
	pugcRouter.Use(throttle())
	//pugcRouter.Use(mustAuth())
	{
		pugcRouter.GET("/addPugc", format(pugc.DefaultPugcController.AddPugc))
		pugcRouter.POST("/ex", format(pugc.DefaultPugcController.ExportExcel))
	}
}
