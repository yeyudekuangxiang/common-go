package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/admin"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util"
)

func adminRouter(router *gin.Engine) {
	adminRouter := router.Group("/admin")
	adminRouter.Use(middleware.AuthAdmin())
	{
		adminRouter.GET("/user", util.Format(admin.DefaultUserController.GetUserInfo))
	}
}
