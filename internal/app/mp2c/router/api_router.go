package router

import (
	"github.com/gin-gonic/gin"
	api2 "mio/internal/app/mp2c/controller/api"
	activity2 "mio/internal/app/mp2c/controller/api/activity"
	authCtr "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/api/coupon"
	"mio/internal/app/mp2c/controller/api/product"
)

func apiRouter(router *gin.Engine) {

	router.GET("/newUser", format(api2.DefaultUserController.GetNewUser))

	//非必须登陆的路由
	authRouter := router.Group("/api/mp2c").Use(auth2(), throttle())
	{

		authRouter.GET("/product-item/list", format(product.DefaultProductController.ProductList))
		authRouter.GET("/openid-coupon/list", format(coupon.DefaultCouponController.CouponListOfOpenid))
		authRouter.POST("/tag/list", format(api2.DefaultTagController.List))
		authRouter.POST("/topic/list", format(api2.DefaultTopicController.List))
		authRouter.GET("/user/get-yzm", format(api2.DefaultUserController.GetYZM))     //获取验证码
		authRouter.GET("/user/check-yzm", format(api2.DefaultUserController.CheckYZM)) //校验验证码

		authRouter.POST("/unidian/callback", api2.DefaultUnidianController.Callback) //手机充值回调函数

		authRouter.POST("/auth/oa/configsign", format(authCtr.DefaultOaController.Sign))

		authRouter.POST("/tool/get-qrcode", format(api2.DefaultToolController.GetQrcode))

		//h5活动页调用
		authRouter.POST("/activity/boc/record", format(activity2.DefaultBocController.FindOrCreateRecord))

		//星星充电订单同步接口
		authRouter.GET("/charge/push", format(api2.DefaultChargeController.Push))

	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c").Use(mustAuth2(), throttle())
	{

		mustAuthRouter.GET("/user", format(api2.DefaultUserController.GetUserInfo))
		mustAuthRouter.GET("/mobile-user", format(api2.DefaultUserController.GetMobileUserInfo))
		mustAuthRouter.GET("/topic/share-qrcode", format(api2.DefaultTopicController.GetShareWeappQrCode))
		mustAuthRouter.POST("/topic/like/change", format(api2.DefaultTopicController.ChangeTopicLike))

		//h5活动页调用
		mustAuthRouter.POST("/activity/boc/answer", format(activity2.DefaultBocController.Answer))
		//小程序端调用
		mustAuthRouter.GET("/activity/boc/share/list", format(activity2.DefaultBocController.GetRecordList))
		mustAuthRouter.GET("/activity/boc/record/mini", format(activity2.DefaultBocController.FindRecordOfMini))
		mustAuthRouter.POST("/activity/bonus/apply", format(activity2.DefaultBocController.ApplySendBonus))

		//GreenMonday活动
		mustAuthRouter.GET("/activity/gm/record", format(activity2.DefaultGMController.GetGMRecord))
		mustAuthRouter.POST("/activity/gm/invitation", format(activity2.DefaultGMController.ReportInvitationRecord))
		mustAuthRouter.POST("/activity/gm/exchange", formatInterface(activity2.DefaultGMController.ExchangeGift))
		mustAuthRouter.POST("/activity/gm/question", format(activity2.DefaultGMController.AnswerQuestion))

		//OCR识别
		mustAuthRouter.POST("/ocr/gm/ticket", format(api2.DefaultOCRController.GmTicket))

		mustAuthRouter.POST("order/submit-from-green", formatInterface(api2.DefaultOrderController.SubmitOrderForGreen))

		mustAuthRouter.GET("duiba/autologin", format(api2.DefaultDuiBaController.AutoLogin))
	}

	openApiRouter := router.Group("/api/mp2c")
	{
		openApiRouter.Any("/duiba/exchange/callback", func(context *gin.Context) {
			context.JSON(200, api2.DefaultDuiBaController.ExchangeCallback(context))
		})

		openApiRouter.Any("/duiba/exchange/result/notice/callback", func(context *gin.Context) {
			context.JSON(200, api2.DefaultDuiBaController.ExchangeCallback(context))
		})
	}

}
