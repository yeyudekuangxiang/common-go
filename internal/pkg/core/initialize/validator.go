package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"gitlab.miotech.com/miotech-application/backend/common-go/validator"
	"log"
)

func InitValidator() {
	log.Println("初始化参数校验组件...")
	binding.Validator = validator.NewValidator()
	log.Println("初始化参数校验组件成功")
}
