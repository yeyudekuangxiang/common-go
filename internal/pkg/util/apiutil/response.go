package apiutil

import (
	"bufio"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"io/ioutil"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/pkg/errno"
	"reflect"
	"strings"
)

func Format(f func(*gin.Context) (gin.H, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		ctx.JSON(FormatErr(err, data))
	}
}

func FormatContent(f func(*gin.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f(ctx)
	}
}

func FormatErr(err error, data interface{}) (int, gin.H) {
	if err != nil && config.Config.App.Debug {
		app.Logger.Debugf("response err:%+v\n", err)
	}

	status, code, message := errno.DecodeErr(err)
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
	return status, FormatResponse(code, data, message)
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
		ctx.JSON(FormatErr(err, data))
	}
}

func FormatCtx(f func(*context.MioContext) (gin.H, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(&context.MioContext{
			Context: ctx,
			DB:      app.DB,
		})
		ctx.JSON(FormatErr(err, data))
	}
}

type AesFormat struct {
	Key []byte
	IV  []byte
}
type AesRequest struct {
	RequestBody string `json:"requestBody" binding:"required"`
}

func (a AesFormat) Format(f func(*gin.Context) (gin.H, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := AesRequest{}
		err := BindForm(ctx, &req)
		if err != nil {
			ctx.JSON(FormatErr(err, nil))
			return
		}
		body, err := encrypttool.AesDecrypt(req.RequestBody, string(a.Key), string(a.IV))
		if err != nil {
			ctx.JSON(FormatErr(err, nil))
			return
		}
		ctx.Request.Body = ioutil.NopCloser(bufio.NewReader(strings.NewReader(body)))

		data, err := f(ctx)
		if err != nil {
			ctx.JSON(FormatErr(err, nil))
			return
		}

		if data == nil {
			data = gin.H{}
		}

		respBody, err := json.Marshal(data)
		if err != nil {
			ctx.JSON(FormatErr(err, nil))
			return
		}
		respBodyStr := encrypttool.AesEncrypt(string(respBody), string(a.Key), string(a.IV))
		ctx.JSON(FormatErr(err, gin.H{
			"responseBody": respBodyStr,
		}))
	}
}
