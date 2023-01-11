package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	activityApi "mio/internal/app/mp2c/controller/api/activity"
	authApi "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/api/badge"
	"mio/internal/app/mp2c/controller/api/business"
	"mio/internal/app/mp2c/controller/api/common"
	"mio/internal/app/mp2c/controller/api/community"
	"mio/internal/app/mp2c/controller/api/coupon"
	"mio/internal/app/mp2c/controller/api/event"
	"mio/internal/app/mp2c/controller/api/message"
	rabbitmqApi "mio/internal/app/mp2c/controller/api/mq"
	"mio/internal/app/mp2c/controller/api/points"
	"mio/internal/app/mp2c/controller/api/product"
	"mio/internal/app/mp2c/controller/api/qnr"
	"mio/internal/app/mp2c/controller/api/question"
	"mio/internal/app/mp2c/controller/open"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
)

func apiRouter(router *gin.Engine) {
	router.GET("/newUser", apiutil.Format(api.DefaultUserController.GetNewUser))
	router.GET("/sendSign", apiutil.Format(message.DefaultMessageController.SendSign))

	//非必须登陆的路由
	authRouter := router.Group("/api/mp2c")
	authRouter.Use(middleware.Auth2(), middleware.Throttle())
	{
		userRouter := authRouter.Group("/user")
		{
			userRouter.GET("/get-yzm", apiutil.Format(api.DefaultUserController.GetYZM))          //获取验证码
			userRouter.GET("/check-yzm", apiutil.Format(api.DefaultUserController.CheckYZM))      //校验验证码
			userRouter.GET("/get-yzm-2b", apiutil.Format(api.DefaultUserController.GetYZM2B))     //获取验证码2b
			userRouter.GET("/check-yzm-2b", apiutil.Format(api.DefaultUserController.CheckYZM2B)) //校验验证码2b
			userRouter.GET("/business/token", apiutil.Format(business.DefaultUserController.GetToken))
		}

		authRouter.GET("/product-item/list", apiutil.Format(product.DefaultProductController.ProductList))
		authRouter.GET("/openid-coupon/list", apiutil.Format(coupon.DefaultCouponController.CouponListOfOpenid))
		//tag
		tagRouter := authRouter.Group("/tag")
		{
			tagRouter.GET("/list", apiutil.Format(community.DefaultTagController.List))
			tagRouter.GET("/detail", apiutil.Format(community.DefaultTagController.DetailTag))
		}

		//社区文章列表
		topicRouter := authRouter.Group("/topic")
		{
			topicRouter.POST("/list", apiutil.Format(community.DefaultTopicController.List))
			topicRouter.GET("/list-topic", apiutil.Format(community.DefaultTopicController.ListTopic))
			topicRouter.GET("/detail", apiutil.Format(community.DefaultTopicController.DetailTopic)) //帖子详情
			topicRouter.GET("/activities/tag-list", apiutil.Format(community.DefaultCommunityActivitiesTagController.List))
		}

		//文章评论列表
		commentRouter := authRouter.Group("/topic/comment")
		{
			commentRouter.GET("/list", apiutil.Format(community.DefaultCommentController.RootList)) //评论列表
			commentRouter.GET("/sub-list", apiutil.Format(community.DefaultCommentController.SubList))
		}

		//统计计数接口
		countRouter := authRouter.Group("/count")
		{
			countRouter.POST("/topic/views")
			countRouter.POST("/topic/collections")
			countRouter.POST("/topic/likes")
			countRouter.POST("/topic/comments")
			countRouter.POST("/comment/likes")
		}

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
		authRouter.GET("upload/sts/token", apiutil.Format(api.DefaultUploadController.GetUploadSTSTokenInfo))
		authRouter.Any("upload/callback", apiutil.Format(api.DefaultUploadController.UploadCallback))
		//星星发券接口限制
		authRouter.POST("/set-exception", apiutil.Format(open.DefaultChargeController.SetException))
		authRouter.POST("/del-exception", apiutil.Format(open.DefaultChargeController.DelException))
		authRouter.GET("icon/list", apiutil.Format(api.DefaultIndexIconController.Page))
	}

	//必须登陆的路由
	mustAuthRouter := router.Group("/api/mp2c")
	mustAuthRouter.Use(middleware.MustAuth2(), middleware.Throttle())
	{
		//答题相关
		qnrRouter := mustAuthRouter.Group("/qnr") //活动答题
		{
			//答题相关路由
			qnrRouter.GET("/subject", apiutil.Format(qnr.DefaultSubjectController.GetList))
			qnrRouter.POST("/create", apiutil.Format(qnr.DefaultSubjectController.Create))
		}
		quizRouter := mustAuthRouter.Group("/quiz") //活动答题
		{
			quizRouter.GET("/daily-questions", apiutil.Format(api.DefaultQuizController.GetDailyQuestions))
			quizRouter.GET("/availability", apiutil.Format(api.DefaultQuizController.Availability))
			quizRouter.POST("/check", apiutil.Format(api.DefaultQuizController.AnswerQuestion))
			quizRouter.POST("/submit", apiutil.Format(api.DefaultQuizController.Submit))
			quizRouter.GET("/daily-result", apiutil.Format(api.DefaultQuizController.DailyResult))
			quizRouter.GET("/summary", apiutil.Format(api.DefaultQuizController.GetSummary))
		}
		questRouter := mustAuthRouter.Group("/question") //通用答题
		{
			//答题相关路由
			questRouter.GET("/subject", apiutil.Format(question.DefaultSubjectController.GetList))
			questRouter.POST("/create", apiutil.Format(question.DefaultSubjectController.Create))
			questRouter.POST("/getUserYearCarbon", apiutil.Format(question.DefaultSubjectController.GetUserYearCarbon))
		}

		//消息相关路由
		messageRouter := mustAuthRouter.Group("/message")
		{
			messageRouter.GET("/sendMessage", apiutil.Format(message.DefaultMessageController.SendMessage))
			messageRouter.GET("/getTemplateId", apiutil.Format(message.DefaultMessageController.GetTemplateId))
			messageRouter.POST("/web-message", apiutil.Format(message.DefaultMessageController.GetWebMessage))
			messageRouter.POST("/web-message-count", apiutil.Format(message.DefaultMessageController.GetWebMessageCount))
			messageRouter.POST("/web-message-haveread", apiutil.Format(message.DefaultMessageController.SetHaveReadWebMessage))
			messageRouter.POST("/im-message-send", apiutil.Format(message.DefaultIMMessageController.Send))
			messageRouter.POST("/im-message-get", apiutil.Format(message.DefaultIMMessageController.GetByFriend))
			messageRouter.POST("/im-message-bind", apiutil.Format(message.DefaultIMMessageController.BindFriend))
		}

		//用户相关路由
		userRouter := mustAuthRouter.Group("/user")
		{
			userRouter.GET("/", apiutil.Format(api.DefaultUserController.GetUserInfo))
			userRouter.GET("/summary", apiutil.Format(api.DefaultUserController.GetUserSummary))
			userRouter.POST("/info/update", apiutil.Format(api.DefaultUserController.UpdateUserInfo))
			userRouter.GET("/account-info", apiutil.Format(api.DefaultUserController.GetUserAccountInfo))
			userRouter.POST("/mobile/bind-by-code", apiutil.Format(api.DefaultUserController.BindMobileByCode))
			userRouter.GET("/mobile/bind-by-yzm", apiutil.Format(api.DefaultUserController.BindMobileByYZM)) //绑定手机
			userRouter.GET("/my-topic", apiutil.Format(community.DefaultTopicController.MyTopic))            //我的帖子列表
			userRouter.GET("/my-reward", apiutil.Format(api.DefaultPointController.MyReward))                //我的奖励
			userRouter.GET("/my-signup", apiutil.Format(community.DefaultTopicController.MySignup))          //我的报名

			userRouter.GET("/topic-collection", apiutil.Format(community.DefaultCollectionController.TopicCollection))    //我的收藏(文章)
			userRouter.POST("/collection", apiutil.Format(community.DefaultCollectionController.Collection))              //收藏(文章)
			userRouter.POST("/cancel-collection", apiutil.Format(community.DefaultCollectionController.CancelCollection)) //取消收藏

			//社区2.0 用户相关
			userRouter.GET("/home-page", apiutil.Format(api.DefaultUserController.HomePage))                      //主页
			userRouter.POST("/update-introduction", apiutil.Format(api.DefaultUserController.UpdateIntroduction)) //更新简介
		}

		//公共接口
		commonRouter := mustAuthRouter.Group("/common")
		{
			commonRouter.GET("/city/list", apiutil.Format(common.DefaultCityController.List))
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

		//社区2.0 文章相关路由
		topicRouter := mustAuthRouter.Group("/topic")
		{
			topicRouter.GET("/share-qrcode", apiutil.Format(community.DefaultTopicController.GetShareWeappQrCode))
			topicRouter.POST("/like/change", apiutil.Format(community.DefaultTopicController.ChangeTopicLike))
			topicRouter.POST("/create", apiutil.Format(community.DefaultTopicController.CreateTopic))
			topicRouter.POST("/update", apiutil.Format(community.DefaultTopicController.UpdateTopic))
			topicRouter.POST("/delete", apiutil.Format(community.DefaultTopicController.DelTopic))
			topicRouter.POST("/activities/signup", apiutil.Format(community.DefaultTopicController.SignupTopic))
			topicRouter.POST("/activities/cancel-signup", apiutil.Format(community.DefaultTopicController.CancelSignupTopic))
			topicRouter.GET("/activities/signup-list", apiutil.Format(community.DefaultTopicController.SignupList))
		}

		//社区2.0 评论相关
		commentRouter := mustAuthRouter.Group("/comment")
		{
			commentRouter.POST("/create", apiutil.Format(community.DefaultCommentController.Create))
			commentRouter.POST("/update", apiutil.Format(community.DefaultCommentController.Update))
			commentRouter.POST("/delete", apiutil.Format(community.DefaultCommentController.Delete))
			commentRouter.POST("/like", apiutil.Format(community.DefaultCommentController.Like))
			commentRouter.GET("/detail", apiutil.Format(community.DefaultCommentController.Detail))
			commentRouter.POST("/turn-comment", apiutil.Format(community.DefaultCommentController.TurnComment))
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
			//checkinRouter.POST("/collect", apiutil.Format(api.DefaultCheckinController.Checkin))
		}

		//扫小票得积分相关路由
		pointCollectRouter := mustAuthRouter.Group("/point-collect")
		{
			pointCollectRouter.POST("/", apiutil.Format(api.DefaultPointCollectController.Collect))
			pointCollectRouter.POST("/new-collect", apiutil.Format(points.DefaultPointsCollectController.ImageCollect))
			pointCollectRouter.POST("/getPageData", apiutil.Format(points.DefaultPointsCollectController.GetPageData))
		}

		//上传文件相关路由
		uploadRouter := mustAuthRouter.Group("/upload")
		{
			uploadRouter.Any("/point-collect", apiutil.Format(api.DefaultUploadController.UploadPointCollectImage))
			uploadRouter.Any("/multipart", apiutil.Format(api.DefaultUploadController.MultipartUploadImage))
			uploadRouter.Any("/", apiutil.Format(api.DefaultUploadController.UploadImage))
		}

		mustAuthRouter.GET("/mobile-user", apiutil.Format(api.DefaultUserController.GetMobileUserInfo))

		//OCR识别
		mustAuthRouter.POST("/ocr/gm/ticket", apiutil.Format(api.DefaultOCRController.GmTicket))

		mustAuthRouter.POST("/order/submit-from-green", apiutil.FormatInterface(api.DefaultOrderController.SubmitOrderForGreen))
		mustAuthRouter.POST("/order/submit-from-event", apiutil.Format(api.DefaultOrderController.SubmitOrderForEvent))
		mustAuthRouter.POST("/order/submit-from-event-gd", apiutil.Format(api.DefaultOrderController.SubmitOrderForEventGD))

		mustAuthRouter.GET("/order/list", apiutil.FormatInterface(api.DefaultOrderController.GetUserOrderList))

		mustAuthRouter.GET("/duiba/autologin", apiutil.Format(api.DefaultDuiBaController.AutoLogin))

		mustAuthRouter.POST("/badge/image", apiutil.Format(badge.DefaultBadgeController.UpdateBadgeImage))
		mustAuthRouter.GET("/badge/list", apiutil.Format(badge.DefaultBadgeController.GetBadgeList))
		mustAuthRouter.POST("/badge/looked", apiutil.Format(badge.DefaultBadgeController.UpdateBadgeIsNew))
		mustAuthRouter.GET("/badge/upload/setting", apiutil.FormatInterface(badge.DefaultBadgeController.UploadOldBadgeImage))

		//兑换券相关
		couponRouter := mustAuthRouter.Group("/coupon")
		{
			couponRouter.GET("/record/list", apiutil.FormatInterface(coupon.DefaultCouponController.GetPageUserCouponRecord))
			couponRouter.POST("/redeem-code", apiutil.FormatInterface(coupon.DefaultCouponController.RedeemCode))
		}

		//获取第三方数据
		platformRouter := mustAuthRouter.Group("/platform")
		{
			platformRouter.GET("/oola-key", apiutil.Format(open.DefaultRecycleController.GetOolaKey))         //获取oolaKey
			platformRouter.POST("/jhx/ticket-create", apiutil.Format(open.DefaultJhxController.TicketCreate)) //金华行-发码
			platformRouter.GET("/jhx/ticket-status", apiutil.Format(open.DefaultJhxController.TicketStatus))  //金华行-查询券码状态
			platformRouter.POST("/all-receive", apiutil.Format(open.DefaultYtxController.AllReceive))         //一键领取
			platformRouter.POST("/pre-point", apiutil.Format(open.DefaultYtxController.PrePointList))         //预加积分列表
		}

		//碳成就相关路由
		carbonRouter := mustAuthRouter.Group("/carbon")
		{
			//	carbonRouter.POST("/create", apiutil.Format(api.DefaultCarbonController.Create))
			//	carbonRouter.POST("/pointToCarbon", apiutil.Format(api.DefaultCarbonController.PointToCarbon))
			carbonRouter.GET("/info", apiutil.Format(api.DefaultCarbonController.Info))
			carbonRouter.GET("/bank", apiutil.Format(api.DefaultCarbonController.Bank))
			carbonRouter.GET("/myBank", apiutil.Format(api.DefaultCarbonController.MyBank))
			carbonRouter.GET("/classify", apiutil.Format(api.DefaultCarbonController.Classify))
			carbonRouter.GET("/history", apiutil.Format(api.DefaultCarbonController.History))
		}
	}

	mqAuthRouter := router.Group("/api/mp2c")
	mqAuthRouter.Use(middleware.MqAuth2(), middleware.Throttle())
	{
		mqRouter := mqAuthRouter.Group("/mq")
		{
			mqRouter.POST("/send_sms", apiutil.Format(rabbitmqApi.DefaultMqController.SendSms))
			mqRouter.POST("/send_zhuge", apiutil.Format(rabbitmqApi.DefaultMqController.SendZhuGe))
			mqRouter.POST("/send_yzm_sms", apiutil.Format(rabbitmqApi.DefaultMqController.SendYzmSms))
		}
	}
}
