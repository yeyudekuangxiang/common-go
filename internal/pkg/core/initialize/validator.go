package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"mio/pkg/validator"
)

func InitValidator() {
	binding.Validator = validator.NewValidator()
}
