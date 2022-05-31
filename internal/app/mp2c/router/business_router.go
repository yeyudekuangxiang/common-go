package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/middleware"
)

func BusinessRouter(router *gin.Engine) {
	businessRouter := router.Group("/api/mp2c/business")
	businessRouter.Use(middleware.Throttle())
	{

	}

	authRouter := router.Group("/api/mp2c/business")
	authRouter.Use(middleware.Throttle())
	{

	}
}
