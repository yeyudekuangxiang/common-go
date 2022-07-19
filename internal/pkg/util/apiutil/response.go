package apiutil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/pkg/errno"
	"mio/pkg/wxwork"
	"os"
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
		case reflect.Ptr:
		default:
			panic("不支持的数据类型")
		}
	}

	if code != 200 && code != errno.ErrAuth.Code() {
		go func() {
			if config.Config.App.Env != "prod" {
				return
			}

			sendErr := wxwork.SendRobotMessage("f0edb1a2-3f9b-4a5d-aa15-9596a32840ec", wxwork.Markdown{
				Content: fmt.Sprintf("**容器:**%s \n\n**来源:**响应 \n\n**消息:**%+v", os.Getenv("HOSTNAME"), err),
			})
			if sendErr != nil {
				log.Printf("推送异常到企业微信失败 %v %v", err, sendErr)
			}
		}()
	}

	return FormatResponse(code, data, message)
}
func FormatResponse(code int, data interface{}, message string) gin.H {
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

func FormatCtx(f func(*context.MioContext) (gin.H, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(&context.MioContext{
			Context: ctx,
			DB:      app.DB,
		})
		ctx.JSON(200, FormatErr(err, data))
	}
}
