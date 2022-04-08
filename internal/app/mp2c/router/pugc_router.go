package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/pugc"
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
