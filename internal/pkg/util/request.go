package util

import (
	"github.com/gin-gonic/gin"
	entity "mio/internal/pkg/model/entity"
	"mio/pkg/errno"
	"mio/pkg/validator"
)

func BindForm(c *gin.Context, data interface{}) error {
	if err := c.ShouldBind(data); err != nil {
		err = validator.TranslateError(err)
		return errno.NewBindErr(err)
	}
	return nil
}
func GetAuthAdmin(c *gin.Context) entity.SystemAdmin {
	if admin, ok := c.Get("AuthAdmin"); ok {
		return admin.(entity.SystemAdmin)
	}
	return entity.SystemAdmin{}
}
func GetAuthUser(c *gin.Context) entity.User {
	if user, ok := c.Get("AuthUser"); ok {
		return user.(entity.User)
	}
	return entity.User{}
}
