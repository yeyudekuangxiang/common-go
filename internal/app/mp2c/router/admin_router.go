package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/admin"
)

func adminRouter(router *gin.Engine) {
	adminRouter := router.Group("/admin")
	adminRouter.Use(authAdmin())
	{
		adminRouter.GET("/user", format(admin.DefaultUserController.GetUserInfo))
	}
}
