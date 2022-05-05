package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/pugc"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
)

func pugcRouter(router *gin.Engine) {
	pugcRouter := router.Group("/pugc")
	pugcRouter.Use(middleware.Throttle())
	//pugcRouter.Use(mustAuth())
	{
		pugcRouter.GET("/addPugc", apiutil.Format(pugc.DefaultPugcController.AddPugc))
		pugcRouter.POST("/ex", apiutil.Format(pugc.DefaultPugcController.ExportExcel))
	}
}
