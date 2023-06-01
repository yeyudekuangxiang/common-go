package middleware

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"io"
	"mio/pkg/errno"
	"net/http"
)

type aesForm struct {
	RequestBody string `json:"requestBody"`
}
type AesResponseWriter struct {
	aesKey []byte
	iv     []byte
	gin.ResponseWriter
}

func (w *AesResponseWriter) Write(data []byte) (int, error) {
	data = encrypttool.AesEncryptCBCPKCS7(data, w.aesKey, w.iv)
	d := base64.StdEncoding.EncodeToString(data)
	return w.ResponseWriter.Write([]byte(fmt.Sprintf(`{"code":200,"data":{"responseBody":"%s"},"message":"OK"}`, d)))
}
func (w *AesResponseWriter) WriteString(str string) (int, error) {
	data := encrypttool.AesEncryptCBCPKCS7([]byte(str), w.aesKey, w.iv)
	d := base64.StdEncoding.EncodeToString(data)
	return w.ResponseWriter.Write([]byte(fmt.Sprintf(`{"code":200,"data":{"responseBody":"%s"},"message":"OK"}`, d)))
}
func decodeErr(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		logx.Errorf("解码失败 %+v %+v", r, err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"code\":%d,\"data\":{},\"message\":\"%s\"}", errno.ErrInternalServer.Code(), errno.ErrInternalServer.Message())))
}
func AesMiddleware(aesKey, iv []byte) gin.HandlerFunc {
	return func(context *gin.Context) {
		data, err := io.ReadAll(context.Request.Body)
		if err != nil {
			decodeErr(context.Writer, context.Request, err)
			return
		}
		_ = context.Request.Body.Close()

		form := aesForm{}
		err = json.Unmarshal(data, &form)
		if err != nil {
			decodeErr(context.Writer, context.Request, err)
			return
		}

		body, err := base64.StdEncoding.DecodeString(form.RequestBody)
		if err != nil {
			decodeErr(context.Writer, context.Request, err)
			return
		}

		body, err = encrypttool.AesDecryptCBCPKCS7(body, aesKey, iv)
		if err != nil {
			decodeErr(context.Writer, context.Request, err)
			return
		}

		context.Request.Body = io.NopCloser(bufio.NewReader(bytes.NewReader(body)))

		w := &AesResponseWriter{
			ResponseWriter: context.Writer,
			aesKey:         aesKey,
			iv:             iv,
		}
		context.Writer = w
	}
}
