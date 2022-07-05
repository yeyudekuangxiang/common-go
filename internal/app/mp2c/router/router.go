package router

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})
	router.GET("/", func(context *gin.Context) {
		context.String(200, "mio")
	})

	router.Any("MP_verify_pp3ZifoA3gboswNR.txt", func(context *gin.Context) {
		context.String(200, "pp3ZifoA3gboswNR")
	})
	router.Any("QUxp4PS6fh.txt", func(context *gin.Context) {
		context.String(200, "c636e427fa1d442771a93ff2885d6c15")
	})

	router.Any("pt04CfOnB5.txt", func(context *gin.Context) {
		context.String(200, "e2082465010d1787e6090c37ed629674")
	})

	apiRouter(router)
	adminRouter(router)
	openRouter(router)
	pugcRouter(router)
}
