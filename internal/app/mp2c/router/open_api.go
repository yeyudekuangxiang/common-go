package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	activityApi "mio/internal/app/mp2c/controller/api/activity"
	authApi "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/open"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
)

func openRouter(router *gin.Engine) {

	openRouter := router.Group("/api/mp2c")
	{
		duibaRouter := openRouter.Group("/duiba")
		{
			duibaRouter.Any("/exchange/callback", func(context *gin.Context) {
				context.JSON(200, api.DefaultDuiBaController.ExchangeCallback(context))
			})

			duibaRouter.Any("/exchange/result/notice/callback", func(context *gin.Context) {
				context.String(200, api.DefaultDuiBaController.ExchangeResultNoticeCallback(context))
			})

			duibaRouter.Any("/order/callback", func(context *gin.Context) {
				context.String(200, api.DefaultDuiBaController.OrderCallback(context))
			})

			duibaRouter.Any("/virtual-good/callback", func(context *gin.Context) {
				context.JSON(200, api.DefaultDuiBaController.VirtualGoodCallback(context))
			})

			duibaRouter.Any("/point/add/callback", func(context *gin.Context) {
				context.JSON(200, api.DefaultDuiBaController.PointAddLogCallback(context))
			})
			duibaRouter.GET("/h5", api.DefaultDuiBaController.DuiBaNoLoginH5)
		}

		oaRouter := openRouter.Group("/oa")
		{
			//微信公众号网页授权登陆
			oaRouter.GET("/auth", func(context *gin.Context) {
				authApi.DefaultOaController.AutoLogin(context)
			})
			//微信公众号网页授权登陆回调
			oaRouter.Any("/auth/callback", func(context *gin.Context) {
				authApi.DefaultOaController.AutoLoginCallback(context)
			})
			//微信公众号网页code登陆
			oaRouter.Any("/login", apiutil.Format(authApi.DefaultOaController.Login))
			//微信网页授权
			oaRouter.POST("/sign", apiutil.Format(authApi.DefaultOaController.Sign))
		}

		openRouter.POST("/weapp/auth", apiutil.Format(authApi.DefaultWeappController.LoginByCode))
		openRouter.GET("/activity/duiba/qr", func(context *gin.Context) {
			if err := activityApi.DefaultZeroController.GetActivityMiniQR(context); err != nil {
				context.String(400, err.Error())
			}
		})
		openRouter.Any("/gitlab/webhook", apiutil.Format(open.DefaultGitlabController.WebHook))
		//绿喵跳转外部平台
		//订单同步接口 （星星充电、快电）
		openRouter.GET("/charge/push", apiutil.Format(open.DefaultChargeController.Push))
		//噢啦\飞蚂蚁旧物回收 订单同步接口
		openRouter.POST("/recycle/oola", apiutil.Format(open.DefaultRecycleController.OolaOrderSync))
		openRouter.POST("/recycle/fmy", apiutil.Format(open.DefaultRecycleController.FmyOrderSync))

		//外部平台调绿喵 需要登陆
		openAuthRouter := openRouter.Group("/auth").Use(middleware.MustAuth2(), middleware.Throttle())
		{
			openAuthRouter.GET("/platform", apiutil.Format(open.DefaultPlatformController.BindPlatformUser))
			openAuthRouter.GET("/platform/bind", apiutil.Format(open.DefaultPlatformController.BindPlatformUser))
		}
		//外部平台调绿喵 不需要登陆
		openRouter.POST("/sync/point", apiutil.Format(open.DefaultPlatformController.SyncPoint))

		openPlatformRouter := openRouter.Group("/platform")
		{
			//金华行
			openBusticketRouter := openPlatformRouter.Group("/busticket")
			{
				openBusticketRouter.POST("/ticket_notify", apiutil.Format(open.DefaultJhxController.BusTicketNotify))  //消费通知
				openBusticketRouter.POST("/get_collect", apiutil.Format(open.DefaultJhxController.GetPreCollectPoint)) //获取积分气泡
				openBusticketRouter.POST("/collect", apiutil.Format(open.DefaultJhxController.CollectPoint))           //收集积分气泡
				openBusticketRouter.POST("/my_account", apiutil.Format(open.DefaultJhxController.MyAccountInfo))
			}
			openPlatformRouter.POST("/pre_collect", apiutil.Format(open.DefaultChargeController.PreCollectPoint))
		}
	}
}
