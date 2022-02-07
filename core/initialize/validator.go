package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"mio/internal/validator"
)

func InitValidator() {
	binding.Validator = validator.NewValidator()
}
