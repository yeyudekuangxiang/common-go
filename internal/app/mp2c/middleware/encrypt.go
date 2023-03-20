package middleware

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"io/ioutil"
	"strings"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Encrypt(key, iv string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		c.Request.Body.Close()

		result, err := encrypttool.AesDecrypt(string(body), key, iv)
		if err != nil {
			panic(err)
		}
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bufio.NewReader(strings.NewReader(result)))
	}
}
