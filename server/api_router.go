package server

import (
	"github.com/gin-gonic/gin"
	"mio/controller/api"
	"mio/controller/api/activity"
	"mio/controller/api/coupon"
	"mio/controller/api/product"
)

func apiRouter(router *gin.Engine) {

	router.GET("/newUser", format(api.DefaultUserController.GetNewUser))

	//非必须登陆的路由
	authRouter := router.Group("/api/mp2c").Use(auth2(), throttle())
	{
		authRouter.GET("/product-item/list", format(product.DefaultProductController.ProductList))
		authRouter.GET("/openid-coupon/list", format(coupon.DefaultCouponController.CouponListOfOpenid))
		authRouter.POST("/tag/list", format(api.DefaultTagController.List))
		authRouter.POST("/topic/list", format(api.DefaultTopicController.List))

		//h5活动页调用
		authRouter.POST("/activity/boc/record", format(activity.DefaultBocController.FindOrCreateRecord))
		//小程序端调用
		authRouter.GET("/activity/boc/share", format(activity.DefaultBocController.GetRecordList))
		authRouter.GET("/activity/boc/record/mini", format(activity.DefaultBocController.FindRecordOfMini))
	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c").Use(mustAuth2(), throttle())
	{
		mustAuthRouter.GET("/user", format(api.DefaultUserController.GetUserInfo))
		mustAuthRouter.GET("/topic/share-qrcode", format(api.DefaultTopicController.GetShareWeappQrCode))
		mustAuthRouter.POST("/topic/like/change", format(api.DefaultTopicController.ChangeTopicLike))

		//h5活动页调用
		mustAuthRouter.POST("/activity/boc/answer", format(activity.DefaultBocController.Answer))
	}

}
