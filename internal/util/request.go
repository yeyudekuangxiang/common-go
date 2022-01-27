package util

import (
	"github.com/gin-gonic/gin"
	"mio/internal/errno"
	"mio/internal/validator"
	"mio/model/entity"
)

func BindForm(c *gin.Context, data interface{}) error {
	if err := c.ShouldBind(data); err != nil {
		err = validator.TranslateError(err)
		return errno.NewBindErr(err)
	}
	return nil
}
func GetAuthAdmin(c *gin.Context) entity.Admin {
	return c.MustGet("AuthAdmin").(entity.Admin)
}
func GetAuthUser(c *gin.Context) entity.User {
	return c.MustGet("AuthUser").(entity.User)
}
