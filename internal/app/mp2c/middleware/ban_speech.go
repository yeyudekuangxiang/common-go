package middleware

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
)

func BanSpeech() gin.HandlerFunc {
	return func(c *gin.Context) {
		usr, exists := c.Get("AuthUser")
		if !exists {
			app.Logger.Error("mustAuth token err")
			c.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrValidation.WithErrMessage("未登录"), nil))
			return
		}
		if usr.(entity.User).Status == 2 {
			app.Logger.Error("validation user status")
			c.AbortWithStatusJSON(apiutil.FormatErr(errno.ErrValidation.WithErrMessage("账号被禁言无法发表评论"), nil))
			return
		}
	}
}
