package server

import (
	"github.com/gin-gonic/gin"
	"mio/controller/api"
	"mio/controller/api/activity"
	authCtr "mio/controller/api/auth"
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
		authRouter.GET("/user/get-yzm", format(api.DefaultUserController.GetYZM))     //获取验证码
		authRouter.GET("/user/check-yzm", format(api.DefaultUserController.CheckYZM)) //校验验证码

		authRouter.POST("/unidian/callback", api.DefaultUnidianController.Callback) //手机充值回调函数

		authRouter.POST("/auth/oa/configsign", format(authCtr.DefaultOaController.Sign))

		authRouter.POST("/tool/get-qrcode", format(api.DefaultToolController.GetQrcode))

		//h5活动页调用
		authRouter.POST("/activity/boc/record", format(activity.DefaultBocController.FindOrCreateRecord))

		//星星充电订单同步接口
		authRouter.GET("/charge/push", format(api.DefaultChargeController.Push))

	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c").Use(mustAuth2(), throttle())
	{

		mustAuthRouter.GET("/user", format(api.DefaultUserController.GetUserInfo))
		mustAuthRouter.GET("/mobile-user", format(api.DefaultUserController.GetMobileUserInfo))
		mustAuthRouter.GET("/topic/share-qrcode", format(api.DefaultTopicController.GetShareWeappQrCode))
		mustAuthRouter.POST("/topic/like/change", format(api.DefaultTopicController.ChangeTopicLike))

		//h5活动页调用
		mustAuthRouter.POST("/activity/boc/answer", format(activity.DefaultBocController.Answer))
		//小程序端调用
		mustAuthRouter.GET("/activity/boc/share/list", format(activity.DefaultBocController.GetRecordList))
		mustAuthRouter.GET("/activity/boc/record/mini", format(activity.DefaultBocController.FindRecordOfMini))
		mustAuthRouter.POST("/activity/bonus/apply", format(activity.DefaultBocController.ApplySendBonus))

		//GreenMonday活动
		mustAuthRouter.GET("/activity/gm/record", format(activity.DefaultGMController.GetGMRecord))
		mustAuthRouter.POST("/activity/gm/invitation", format(activity.DefaultGMController.ReportInvitationRecord))
		mustAuthRouter.POST("/activity/gm/exchange", formatInterface(activity.DefaultGMController.ExchangeGift))
		mustAuthRouter.POST("/activity/gm/question", format(activity.DefaultGMController.AnswerQuestion))

		//OCR识别
		mustAuthRouter.POST("/ocr/gm/ticket", format(api.DefaultOCRController.GmTicket))

		mustAuthRouter.POST("order/submit-from-green", formatInterface(api.DefaultOrderController.SubmitOrderForGreen))

		mustAuthRouter.GET("duiba/autologin", format(api.DefaultDuiBaController.AutoLogin))
	}

	openApiRouter := router.Group("/api/mp2c")
	{
		openApiRouter.Any("/duiba/exchange/callback", func(context *gin.Context) {
			context.JSON(200, api.DefaultDuiBaController.ExchangeCallback(context))
		})

		openApiRouter.Any("/duiba/exchange/result/notice/callback", func(context *gin.Context) {
			context.JSON(200, api.DefaultDuiBaController.ExchangeCallback(context))
		})
	}

}
