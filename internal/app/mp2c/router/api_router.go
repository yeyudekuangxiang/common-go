package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	activityApi "mio/internal/app/mp2c/controller/api/activity"
	authApi "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/api/badge"
	"mio/internal/app/mp2c/controller/api/business"
	"mio/internal/app/mp2c/controller/api/coupon"
	"mio/internal/app/mp2c/controller/api/event"
	"mio/internal/app/mp2c/controller/api/product"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
)

func apiRouter(router *gin.Engine) {

	router.GET("/newUser", apiutil.Format(api.DefaultUserController.GetNewUser))

	//非必须登陆的路由
	authRouter := router.Group("/api/mp2c")
	authRouter.Use(middleware.Auth2(), middleware.Throttle())
	{
		userRouter := authRouter.Group("/user")
		{
			userRouter.GET("/get-yzm", apiutil.Format(api.DefaultUserController.GetYZM))     //获取验证码
			userRouter.GET("/check-yzm", apiutil.Format(api.DefaultUserController.CheckYZM)) //校验验证码
			//企业用户获取token
			userRouter.GET("/business/token", apiutil.Format(business.DefaultUserController.GetToken))
		}

		authRouter.GET("/product-item/list", apiutil.Format(product.DefaultProductController.ProductList))
		authRouter.GET("/openid-coupon/list", apiutil.Format(coupon.DefaultCouponController.CouponListOfOpenid))
		authRouter.POST("/tag/list", apiutil.Format(api.DefaultTagController.List))
		authRouter.POST("/topic/list", apiutil.Format(api.DefaultTopicController.List))

		authRouter.POST("/unidian/callback", api.DefaultUnidianController.Callback) //手机充值回调函数

		authRouter.POST("/auth/oa/configsign", apiutil.Format(authApi.DefaultOaController.Sign))

		authRouter.POST("/tool/get-qrcode", apiutil.Format(api.DefaultToolController.GetQrcode))

		//h5活动页调用
		authRouter.POST("/activity/boc/record", apiutil.Format(activityApi.DefaultBocController.FindOrCreateRecord))
		//广东小学图书馆公益捐书活动
		authRouter.POST("/activity/answer/homepage", apiutil.Format(activityApi.DefaultAnswerController.HomePage))

		eventRouter := authRouter.Group("/event")
		{
			eventRouter.GET("category/list", apiutil.Format(event.DefaultEventController.GetEventCategoryList))
			eventRouter.GET("/list", apiutil.Format(event.DefaultEventController.GetEventList))
			eventRouter.GET("detail", apiutil.Format(event.DefaultEventController.GetEventFullDetail))
		}

		authRouter.GET("banner/list", apiutil.Format(api.DefaultBannerController.GetBannerList))
		authRouter.GET("upload/token", apiutil.Format(api.DefaultUploadController.GetUploadTokenInfo))
		authRouter.GET("upload/callback", apiutil.Format(api.DefaultUploadController.UploadCallback))
	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c")
	mustAuthRouter.Use(middleware.MustAuth2(), middleware.Throttle())
	{
		//用户相关路由
		userRouter := mustAuthRouter.Group("/user")
		{
			userRouter.GET("/", apiutil.Format(api.DefaultUserController.GetUserInfo))
			userRouter.POST("/mobile/bind-by-code", apiutil.Format(api.DefaultUserController.BindMobileByCode))
			userRouter.GET("/summary", apiutil.Format(api.DefaultUserController.GetUserSummary))
			userRouter.POST("/info/update", apiutil.Format(api.DefaultUserController.UpdateUserInfo))
			userRouter.GET("/account-info", apiutil.Format(api.DefaultUserController.GetUserAccountInfo))
		}
		//邀请得积分
		inviteRouter := mustAuthRouter.Group("/invite")
		{
			inviteRouter.GET("/qrcode", apiutil.Format(api.DefaultInviteController.GetShareQrCode))
			inviteRouter.GET("/list", apiutil.Format(api.DefaultInviteController.GetInviteList))
		}

		//活动相关路由
		activityRouter := mustAuthRouter.Group("/activity")
		{
			//h5活动页调用
			activityRouter.POST("/boc/answer", apiutil.Format(activityApi.DefaultBocController.Answer))
			//小程序端调用
			activityRouter.GET("/boc/share/list", apiutil.Format(activityApi.DefaultBocController.GetRecordList))
			activityRouter.GET("/boc/record/mini", apiutil.Format(activityApi.DefaultBocController.FindRecordOfMini))
			activityRouter.POST("/bonus/apply", apiutil.Format(activityApi.DefaultBocController.ApplySendBonus))

			//GreenMonday活动
			activityRouter.GET("/gm/record", apiutil.Format(activityApi.DefaultGMController.GetGMRecord))
			activityRouter.POST("/gm/invitation", apiutil.Format(activityApi.DefaultGMController.ReportInvitationRecord))
			activityRouter.POST("/gm/exchange", apiutil.FormatInterface(activityApi.DefaultGMController.ExchangeGift))
			activityRouter.POST("/gm/question", apiutil.Format(activityApi.DefaultGMController.AnswerQuestion))
			activityRouter.POST("/zero/autologin", apiutil.Format(activityApi.DefaultZeroController.AutoLogin))
			activityRouter.POST("/zero/storeurl", apiutil.Format(activityApi.DefaultZeroController.StoreUrl))
			activityRouter.POST("/duiba/autologin", apiutil.Format(activityApi.DefaultZeroController.DuiBaAutoLogin))
			activityRouter.POST("/duiba/storeurl", apiutil.Format(activityApi.DefaultZeroController.DuiBaStoreUrl))
			//广东小学图书馆公益捐书活动
			activityRouter.POST("/answer/start-question", apiutil.Format(activityApi.DefaultAnswerController.StartQuestion))
			activityRouter.POST("/answer/end-question", apiutil.Format(activityApi.DefaultAnswerController.EndQuestion))
			activityRouter.POST("/answer/create-school", apiutil.Format(activityApi.DefaultAnswerController.CreateSchool))
			activityRouter.GET("/answer/get-city-list", apiutil.Format(activityApi.DefaultAnswerController.GetCityList))
			activityRouter.GET("/answer/get-grade-list", apiutil.Format(activityApi.DefaultAnswerController.GetGradeList))
			activityRouter.GET("/answer/get-school-list", apiutil.Format(activityApi.DefaultAnswerController.GetSchoolList))
			activityRouter.GET("/answer/get-achievement", apiutil.Format(activityApi.DefaultAnswerController.GetAchievement))
			activityRouter.POST("/answer/put-file", apiutil.Format(activityApi.DefaultAnswerController.PutFile))
			activityRouter.GET("/answer/get-user-school", apiutil.Format(activityApi.DefaultAnswerController.GetUserSchool))
			activityRouter.POST("/answer/close-late-tips", apiutil.Format(activityApi.DefaultAnswerController.CloseLateTips))
		}

		//酷喵圈相关路由
		topicRouter := mustAuthRouter.Group("/topic")
		{
			topicRouter.GET("/share-qrcode", apiutil.Format(api.DefaultTopicController.GetShareWeappQrCode))
			topicRouter.POST("/like/change", apiutil.Format(api.DefaultTopicController.ChangeTopicLike))
			//topicRouter.POST("/create", apiutil.Format(api.DefaultTopicController.CreateTopic))
		}

		//积分相关路由
		pointRouter := mustAuthRouter.Group("/point")
		{
			pointRouter.Any("/list", apiutil.Format(api.DefaultPointController.GetPointTransactionList))
			pointRouter.GET("/", apiutil.Format(api.DefaultPointController.GetPoint))
		}

		//步行相关的路由
		stepRouter := mustAuthRouter.Group("/step")
		{
			//更新用户步数历史记录
			stepRouter.POST("/update", apiutil.Format(api.DefaultStepController.UpdateStepTotal))
			//获取最近7天内步行数据
			stepRouter.GET("/weekly-history", apiutil.FormatInterface(api.DefaultStepController.WeeklyHistory))
			stepRouter.POST("/collect", apiutil.Format(api.DefaultStepController.Collect))
		}

		//签到相关路由
		checkinRouter := mustAuthRouter.Group("/checkin")
		{
			checkinRouter.GET("/info", apiutil.Format(api.DefaultCheckinController.GetCheckinInfo))
			checkinRouter.POST("/collect", apiutil.Format(api.DefaultCheckinController.Checkin))
		}

		//答题相关路由
		quizRouter := mustAuthRouter.Group("/quiz")
		{
			quizRouter.GET("/daily-questions", apiutil.Format(api.DefaultQuizController.GetDailyQuestions))
			quizRouter.GET("/availability", apiutil.Format(api.DefaultQuizController.Availability))
			quizRouter.POST("/check", apiutil.Format(api.DefaultQuizController.AnswerQuestion))
			quizRouter.POST("/submit", apiutil.Format(api.DefaultQuizController.Submit))
			quizRouter.GET("/daily-result", apiutil.Format(api.DefaultQuizController.DailyResult))
			quizRouter.GET("/summary", apiutil.Format(api.DefaultQuizController.GetSummary))
		}

		//扫小票得积分相关路由
		pointCollectRouter := mustAuthRouter.Group("/point-collect")
		{
			pointCollectRouter.POST("", apiutil.Format(api.DefaultPointCollectController.Collect))
		}

		//上传文件相关路由
		uploadRouter := mustAuthRouter.Group("/upload")
		{
			uploadRouter.Any("/point-collect", apiutil.Format(api.DefaultUploadController.UploadPointCollectImage))
		}

		mustAuthRouter.GET("/mobile-user", apiutil.Format(api.DefaultUserController.GetMobileUserInfo))

		//OCR识别
		mustAuthRouter.POST("/ocr/gm/ticket", apiutil.Format(api.DefaultOCRController.GmTicket))

		mustAuthRouter.POST("/order/submit-from-green", apiutil.FormatInterface(api.DefaultOrderController.SubmitOrderForGreen))
		mustAuthRouter.POST("/order/submit-from-event", apiutil.Format(api.DefaultOrderController.SubmitOrderForEvent))
		mustAuthRouter.GET("/order/list", apiutil.FormatInterface(api.DefaultOrderController.GetUserOrderList))

		mustAuthRouter.GET("/duiba/autologin", apiutil.Format(api.DefaultDuiBaController.AutoLogin))

		mustAuthRouter.POST("/badge/image", apiutil.Format(badge.DefaultBadgeController.UpdateBadgeImage))
		mustAuthRouter.GET("/badge/list", apiutil.Format(badge.DefaultBadgeController.GetBadgeList))
		mustAuthRouter.POST("/badge/looked", apiutil.Format(badge.DefaultBadgeController.UpdateBadgeIsNew))

		//兑换券相关
		mustAuthRouter.GET("/coupon/record/list", apiutil.FormatInterface(coupon.DefaultCouponController.GetPageUserCouponRecord))

		//碳成就相关路由
		carbonRouter := mustAuthRouter.Group("/carbon")
		{
			carbonRouter.POST("/create", apiutil.Format(api.DefaultCarbonController.Create))
			carbonRouter.GET("/info", apiutil.Format(api.DefaultCarbonController.Info))
			carbonRouter.GET("/bank", apiutil.Format(api.DefaultCarbonController.Bank))
			carbonRouter.GET("/myBank", apiutil.Format(api.DefaultCarbonController.MyBank))
			carbonRouter.GET("/classify", apiutil.Format(api.DefaultCarbonController.Classify))
			carbonRouter.GET("/history", apiutil.Format(api.DefaultCarbonController.History))
			carbonRouter.GET("/addClassify", apiutil.Format(api.DefaultCarbonController.AddClassify))
			carbonRouter.GET("/addHistory", apiutil.Format(api.DefaultCarbonController.AddHistory))
		}

	}

}
