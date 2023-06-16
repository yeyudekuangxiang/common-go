package router

import (
	"mio/internal/app/mp2c/controller/api"
	activityApi "mio/internal/app/mp2c/controller/api/activity"
	authApi "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/open"
	"mio/internal/app/mp2c/middleware"
	"mio/internal/pkg/util/apiutil"
	"net/http"

	"github.com/gin-gonic/gin"
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

		openRouter.GET("zyh", apiutil.Format(open.DefaultZyhController.Zyh))
		openRouter.POST("sendpoint", apiutil.Format(open.DefaultZyhController.SendPoint))

		oaRouter := openRouter.Group("/oa")
		{

			//微信公众号网页授权登陆
			oaRouter.GET("/auth", func(context *gin.Context) {
				//登录改造，公众号授权现在没有改造且现在没用到，先关闭此接口
				context.String(http.StatusOK, "系统开小差了,请稍后再试")
				return
				authApi.DefaultOaController.AutoLogin(context)
			})
			//微信公众号网页授权登陆回调
			oaRouter.Any("/auth/callback", func(context *gin.Context) {
				//登录改造，公众号授权现在没有改造且现在没用到，先关闭此接口
				context.String(http.StatusOK, "系统开小差了,请稍后再试")
				return
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

		carbonRouth := openRouter.Group("/carbon")
		{
			carbonRouth.POST("/mio-business/change/callback", func(context *gin.Context) {
				res, err := api.DefaultCarbonController.Create(context)
				if err != nil {
					context.String(500, err.Error())
					return
				}
				context.JSON(200, res)
			})
		}

		//外部平台调绿喵 需要登陆
		openAuthRouter := openRouter.Group("/auth").Use(middleware.MustAuth(), middleware.Throttle())
		{
			openAuthRouter.GET("/platform", apiutil.Format(open.DefaultPlatformController.BindPlatformUser))
			openAuthRouter.GET("/platform/bind", apiutil.Format(open.DefaultPlatformController.BindPlatformUser))
			openAuthRouter.POST("/check/msg", apiutil.Format(open.DefaultPlatformController.CheckMgs))
			openAuthRouter.POST("/check/media", apiutil.Format(open.DefaultPlatformController.CheckMedia))
		}
		//外部平台调绿喵 不需要登陆
		openRouter.POST("/sync/point", apiutil.Format(open.DefaultPlatformController.SyncPoint))

		openPlatformRouter := openRouter.Group("/platform")
		{
			//金华行
			openBusticketRouter := openPlatformRouter.Group("/busticket")
			{
				openBusticketRouter.POST("/ticket_notify", apiutil.Format(open.DefaultJhxController.BusTicketNotify))     //消费通知
				openBusticketRouter.POST("/get_collect", apiutil.Format(open.DefaultJhxController.JhxGetPreCollectPoint)) //获取积分气泡
				openBusticketRouter.POST("/collect", apiutil.Format(open.DefaultJhxController.JhxCollectPoint))           //收集积分气泡
			}
			openPlatformRouter.POST("/pre_collect", apiutil.Format(open.DefaultJhxController.JhxPreCollectPoint))         //金华行单独调用 预加积分
			openPlatformRouter.POST("/pre_point", apiutil.Format(open.DefaultPlatformController.PrePoint))                //第三方平台 预加积分
			openPlatformRouter.POST("/get_pre_point", apiutil.Format(open.DefaultPlatformController.GetPrePointList))     //第三方平台 获取积分
			openPlatformRouter.POST("/collect_pre_point", apiutil.Format(open.DefaultPlatformController.CollectPrePoint)) //第三方平台 消费积分
			openPlatformRouter.POST("/syncusr", apiutil.Format(open.DefaultPlatformController.Syncusr))                   //注册回调
			openPlatformRouter.POST("/recycle", apiutil.Format(open.DefaultRecycleController.Recycle))                    //旧物回收
			//
			openPlatformRouter.POST("/charge/ykc", apiutil.Format(open.DefaultChargeController.Ykc))
		}

	}
}
