package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"mio/config"
	"mio/pkg/errno"
	"reflect"
)

func formatErr(err error, data interface{}) gin.H {
	code, message := errno.DecodeErr(err)
	if data == nil {
		data = make(map[string]interface{})
	} else {
		val := reflect.ValueOf(data)
		switch val.Kind() {
		case reflect.Map, reflect.Slice, reflect.Interface:
			if val.IsNil() {
				data = make(map[string]interface{})
			}
		case reflect.Struct:
		default:
			panic("不支持的数据类型")
		}
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
func formatInterface(f func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		if err != nil && config.Config.App.Debug {
			log.Printf("%+v\n", err)
		}
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
	router.Any("QUxp4PS6fh.txt", func(context *gin.Context) {
		context.String(200, "c636e427fa1d442771a93ff2885d6c15")
	})

	apiRouter(router)
	adminRouter(router)
	pugcRouter(router)
}
