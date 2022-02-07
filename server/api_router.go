package server

import (
	"github.com/gin-gonic/gin"
	"mio/controller/api"
)

func apiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(throttle())
	{
		apiRouter.GET("/user", format(api.DefaultUserController.GetUserInfo))
		apiRouter.GET("/newUser", format(api.DefaultUserController.GetNewUser))
	}
}
