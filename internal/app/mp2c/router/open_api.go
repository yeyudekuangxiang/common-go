package router

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api"
	activityApi "mio/internal/app/mp2c/controller/api/activity"
	authApi "mio/internal/app/mp2c/controller/api/auth"
	"mio/internal/app/mp2c/controller/api/system"
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
		openRouter.Any("/gitlab/callback", apiutil.Format(system.DefaultGitlabController.Callback))
		openRouter.GET("/activity/duiba/qr", func(context *gin.Context) {
			if err := activityApi.DefaultZeroController.GetActivityMiniQR(context); err != nil {
				context.String(400, err.Error())
			}
		})
	}
}
