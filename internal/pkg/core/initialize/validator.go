package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"log"
	"mio/pkg/validator"
)

func InitValidator() {
	log.Println("初始化参数校验组件...")
	binding.Validator = validator.NewValidator()
	log.Println("初始化参数校验组件成功")
}
