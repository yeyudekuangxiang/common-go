package apitool

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/yeyudekuangxiang/common-go/validator"
	"net/http"
	"strings"
)

func SetValidator(structValidator binding.StructValidator) {
	binding.Validator = structValidator
}

func BindForm(req *http.Request, data interface{}) error {
	if err := ShouldBind(req, data); err != nil {
		err = validator.TranslateError(err)
		return err
	}
	return nil
}
func ShouldBind(req *http.Request, data interface{}) error {
	b := binding.Default(req.Method, getContentType(req))
	return ShouldBindWith(req, data, b)
}

func ShouldBindWith(req *http.Request, obj interface{}, b binding.Binding) error {
	return b.Bind(req, obj)
}

func getContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if contentType != "" {
		contentTypeList := strings.Split(contentType, ";")
		for _, c := range contentTypeList {
			c := strings.TrimSpace(c)
			if strings.HasPrefix(c, "application/") {
				contentType = c
				break
			}
		}
	}
	return contentType
}
