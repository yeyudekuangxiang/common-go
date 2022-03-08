package server

import (
	"github.com/gin-gonic/gin"
	"mio/controller/api"
	"mio/controller/api/activity"
	"mio/controller/api/coupon"
	"mio/controller/api/product"
)

func apiRouter(router *gin.Engine) {

	//非必须登陆的路由
	authRouter := router.Group("/api").Use(auth2(), throttle())
	{
		authRouter.GET("/newUser", format(api.DefaultUserController.GetNewUser))
		authRouter.GET("/mp2c/product-item/list", format(product.DefaultProductController.ProductList))
		authRouter.GET("/mp2c/openid-coupon/list", format(coupon.DefaultCouponController.CouponListOfOpenid))
		authRouter.POST("/mp2c/tag/list", format(api.DefaultTagController.List))
		authRouter.POST("/mp2c/topic/list", format(api.DefaultTopicController.List))
	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api").Use(mustAuth2(), throttle())
	{
		mustAuthRouter.GET("/mp2c/user", format(api.DefaultUserController.GetUserInfo))
		mustAuthRouter.GET("/mp2c/topic/share-qrcode", format(api.DefaultTopicController.GetShareWeappQrCode))
		mustAuthRouter.POST("/mp2c/topic/like/change", format(api.DefaultTopicController.ChangeTopicLike))

		mustAuthRouter.GET("/activity/boc/apply/share", format(activity.DefaultBocController.GetRecordList))
		mustAuthRouter.POST("/activity/boc/apply/add", format(activity.DefaultBocController.AddRecord))
	}

}
