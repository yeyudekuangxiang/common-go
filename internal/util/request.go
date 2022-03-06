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
	if admin, ok := c.Get("AuthAdmin"); ok {
		return admin.(entity.Admin)
	}
	return entity.Admin{}
}
func GetAuthUser(c *gin.Context) entity.User {
	if user, ok := c.Get("AuthUser"); ok {
		return user.(entity.User)
	}
	return entity.User{}
}
