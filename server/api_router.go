package server

import (
	"github.com/gin-gonic/gin"
	"mio/controller/api"
	"mio/controller/api/coupon"
	"mio/controller/api/product"
)

func apiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(throttle())
	{
		apiRouter.GET("/user", format(api.DefaultUserController.GetUserInfo))
		apiRouter.GET("/newUser", format(api.DefaultUserController.GetNewUser))
		apiRouter.GET("/mp2c/product-item/list", format(product.DefaultProductController.ProductList))
		apiRouter.GET("/mp2c/openid-coupon/list", format(coupon.DefaultCouponController.CouponListOfOpenid))
		apiRouter.POST("/mp2c/topic/list", format(api.DefaultTopicController.List))
		apiRouter.POST("/mp2c/tag/list", format(api.DefaultTagController.List))

	}

	authRouter := router.Group("/api").Use(mustAuth2())
	{
		authRouter.GET("/topic/share-qrcode", format(api.DefaultTopicController.GetShareWeappQrCode))
		authRouter.POST("/topic/like/change", format(api.DefaultTopicController.ChangeTopicLike))
	}
}
