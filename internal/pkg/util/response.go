package util

import (
	"github.com/gin-gonic/gin"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/errno"
	"reflect"
)

func Format(f func(*gin.Context) (gin.H, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		ctx.JSON(200, FormatErr(err, data))
	}
}
func FormatErr(err error, data interface{}) gin.H {
	if err != nil && config.Config.App.Debug {
		log.Printf("response  err:%+v\n", err)
		app.Logger.Debugf("response err:%+v\n", err)
	}

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
func FormatInterface(f func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		ctx.JSON(200, FormatErr(err, data))
	}
}
