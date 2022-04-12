package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	activityApi "mio/internal/app/mp2c/controller/api/activity"
	authApi "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/api/coupon"
	"mio/internal/app/mp2c/controller/api/product"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util"
)

func apiRouter(router *gin.Engine) {

	router.GET("/newUser", util.Format(api.DefaultUserController.GetNewUser))

	//非必须登陆的路由
	authRouter := router.Group("/api/mp2c").Use(middleware.Auth2(), middleware.Throttle())
	{

		authRouter.GET("/product-item/list", util.Format(product.DefaultProductController.ProductList))
		authRouter.GET("/openid-coupon/list", util.Format(coupon.DefaultCouponController.CouponListOfOpenid))
		authRouter.POST("/tag/list", util.Format(api.DefaultTagController.List))
		authRouter.POST("/topic/list", util.Format(api.DefaultTopicController.List))
		authRouter.GET("/user/get-yzm", util.Format(api.DefaultUserController.GetYZM))     //获取验证码
		authRouter.GET("/user/check-yzm", util.Format(api.DefaultUserController.CheckYZM)) //校验验证码

		authRouter.POST("/unidian/callback", api.DefaultUnidianController.Callback) //手机充值回调函数

		authRouter.POST("/auth/oa/configsign", util.Format(authApi.DefaultOaController.Sign))

		authRouter.POST("/tool/get-qrcode", util.Format(api.DefaultToolController.GetQrcode))

		//h5活动页调用
		authRouter.POST("/activity/boc/record", util.Format(activityApi.DefaultBocController.FindOrCreateRecord))

		//星星充电订单同步接口
		authRouter.GET("/charge/push", util.Format(api.DefaultChargeController.Push))

	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c").Use(middleware.MustAuth2(), middleware.Throttle())
	{

		mustAuthRouter.GET("/user", util.Format(api.DefaultUserController.GetUserInfo))
		mustAuthRouter.GET("/mobile-user", util.Format(api.DefaultUserController.GetMobileUserInfo))
		mustAuthRouter.GET("/topic/share-qrcode", util.Format(api.DefaultTopicController.GetShareWeappQrCode))
		mustAuthRouter.POST("/topic/like/change", util.Format(api.DefaultTopicController.ChangeTopicLike))

		//h5活动页调用
		mustAuthRouter.POST("/activity/boc/answer", util.Format(activityApi.DefaultBocController.Answer))
		//小程序端调用
		mustAuthRouter.GET("/activity/boc/share/list", util.Format(activityApi.DefaultBocController.GetRecordList))
		mustAuthRouter.GET("/activity/boc/record/mini", util.Format(activityApi.DefaultBocController.FindRecordOfMini))
		mustAuthRouter.POST("/activity/bonus/apply", util.Format(activityApi.DefaultBocController.ApplySendBonus))

		//GreenMonday活动
		mustAuthRouter.GET("/activity/gm/record", util.Format(activityApi.DefaultGMController.GetGMRecord))
		mustAuthRouter.POST("/activity/gm/invitation", util.Format(activityApi.DefaultGMController.ReportInvitationRecord))
		mustAuthRouter.POST("/activity/gm/exchange", util.FormatInterface(activityApi.DefaultGMController.ExchangeGift))
		mustAuthRouter.POST("/activity/gm/question", util.Format(activityApi.DefaultGMController.AnswerQuestion))
		mustAuthRouter.POST("/activity/zero/autologin", util.Format(activityApi.DefaultZeroController.AutoLogin))
		mustAuthRouter.POST("/activity/zero/storeurl", util.Format(activityApi.DefaultZeroController.StoreUrl))

		//OCR识别
		mustAuthRouter.POST("/ocr/gm/ticket", util.Format(api.DefaultOCRController.GmTicket))

		mustAuthRouter.POST("order/submit-from-green", util.FormatInterface(api.DefaultOrderController.SubmitOrderForGreen))

		mustAuthRouter.GET("duiba/autologin", util.Format(api.DefaultDuiBaController.AutoLogin))
	}

	openApiRouter := router.Group("/api/mp2c")
	{
		openApiRouter.Any("/duiba/exchange/callback", func(context *gin.Context) {
			context.JSON(200, api.DefaultDuiBaController.ExchangeCallback(context))
		})

		openApiRouter.Any("/duiba/exchange/result/notice/callback", func(context *gin.Context) {
			context.JSON(200, api.DefaultDuiBaController.ExchangeCallback(context))
		})

		//微信公众号网页授权登陆
		openApiRouter.GET("/oa/auth", func(context *gin.Context) {
			authApi.DefaultOaController.AutoLogin(context)
		})
		//微信公众号网页授权登陆回调
		openApiRouter.Any("/oa/auth/callback", func(context *gin.Context) {
			authApi.DefaultOaController.AutoLoginCallback(context)
		})
		//微信公众号网页code登陆
		openApiRouter.Any("/oa/login", util.Format(authApi.DefaultOaController.Login))
		//微信网页授权
		openApiRouter.POST("/oa/sign", util.Format(authApi.DefaultOaController.Sign))
	}

}
