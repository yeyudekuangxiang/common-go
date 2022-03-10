package server

import (
	"github.com/gin-gonic/gin"
	"mio/internal/errno"
	"reflect"
)

func formatErr(err error, data interface{}) gin.H {
	code, message := errno.DecodeErr(err)
	if data == nil || reflect.ValueOf(data).IsNil() {
		data = make(map[string]interface{})
	}
	return gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
}
func format(f func(*gin.Context) (gin.H, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		ctx.JSON(200, formatErr(err, data))
	}
}
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

	apiRouter(router)
	adminRouter(router)
	pugcRouter(router)
}
