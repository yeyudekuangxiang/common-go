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
	authRouter := router.Group("/api/mp2c")
	authRouter.Use(middleware.Auth2(), middleware.Throttle())
	{

		userRouter := authRouter.Group("/user")
		{
			userRouter.GET("/get-yzm", util.Format(api.DefaultUserController.GetYZM))     //获取验证码
			userRouter.GET("/check-yzm", util.Format(api.DefaultUserController.CheckYZM)) //校验验证码
		}

		authRouter.GET("/product-item/list", util.Format(product.DefaultProductController.ProductList))
		authRouter.GET("/openid-coupon/list", util.Format(coupon.DefaultCouponController.CouponListOfOpenid))
		authRouter.POST("/tag/list", util.Format(api.DefaultTagController.List))
		authRouter.POST("/topic/list", util.Format(api.DefaultTopicController.List))

		authRouter.POST("/unidian/callback", api.DefaultUnidianController.Callback) //手机充值回调函数

		authRouter.POST("/auth/oa/configsign", util.Format(authApi.DefaultOaController.Sign))

		authRouter.POST("/tool/get-qrcode", util.Format(api.DefaultToolController.GetQrcode))

		//h5活动页调用
		authRouter.POST("/activity/boc/record", util.Format(activityApi.DefaultBocController.FindOrCreateRecord))

		//星星充电订单同步接口
		authRouter.GET("/charge/push", util.Format(api.DefaultChargeController.Push))

	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c")
	mustAuthRouter.Use(middleware.MustAuth2(), middleware.Throttle())
	{
		//用户相关路由
		userRouter := mustAuthRouter.Group("/user")
		{
			userRouter.GET("/", util.Format(api.DefaultUserController.GetUserInfo))
			userRouter.POST("/mobile/bind-by-code", util.Format(api.DefaultUserController.BindMobileByCode))
			userRouter.GET("/summary", util.Format(api.DefaultUserController.GetUserSummary))
		}

		//活动相关路由
		activityRouter := mustAuthRouter.Group("/activity")
		{
			//h5活动页调用
			activityRouter.POST("/boc/answer", util.Format(activityApi.DefaultBocController.Answer))
			//小程序端调用
			activityRouter.GET("/boc/share/list", util.Format(activityApi.DefaultBocController.GetRecordList))
			activityRouter.GET("/boc/record/mini", util.Format(activityApi.DefaultBocController.FindRecordOfMini))
			activityRouter.POST("/bonus/apply", util.Format(activityApi.DefaultBocController.ApplySendBonus))

			//GreenMonday活动
			activityRouter.GET("/gm/record", util.Format(activityApi.DefaultGMController.GetGMRecord))
			activityRouter.POST("/gm/invitation", util.Format(activityApi.DefaultGMController.ReportInvitationRecord))
			activityRouter.POST("/gm/exchange", util.FormatInterface(activityApi.DefaultGMController.ExchangeGift))
			activityRouter.POST("/gm/question", util.Format(activityApi.DefaultGMController.AnswerQuestion))
			activityRouter.POST("/zero/autologin", util.Format(activityApi.DefaultZeroController.AutoLogin))
			activityRouter.POST("/zero/storeurl", util.Format(activityApi.DefaultZeroController.StoreUrl))
			activityRouter.POST("/duiba/autologin", util.Format(activityApi.DefaultZeroController.DuiBaAutoLogin))
			activityRouter.POST("/duiba/storeurl", util.Format(activityApi.DefaultZeroController.DuiBaStoreUrl))
		}

		//酷喵圈相关路由
		topicRouter := mustAuthRouter.Group("/topic")
		{
			topicRouter.GET("/share-qrcode", util.Format(api.DefaultTopicController.GetShareWeappQrCode))
			topicRouter.POST("/like/change", util.Format(api.DefaultTopicController.ChangeTopicLike))
		}

		//积分相关路由
		pointRouter := mustAuthRouter.Group("/point")
		{
			pointRouter.Any("/list", util.Format(api.DefaultPointController.GetPointTransactionList))
			pointRouter.GET("/", util.Format(api.DefaultPointController.GetPoint))
		}

		mustAuthRouter.GET("/mobile-user", util.Format(api.DefaultUserController.GetMobileUserInfo))

		//OCR识别
		mustAuthRouter.POST("/ocr/gm/ticket", util.Format(api.DefaultOCRController.GmTicket))

		mustAuthRouter.POST("/order/submit-from-green", util.FormatInterface(api.DefaultOrderController.SubmitOrderForGreen))

		mustAuthRouter.GET("/duiba/autologin", util.Format(api.DefaultDuiBaController.AutoLogin))

	}

}
