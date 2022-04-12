package util

import (
	"github.com/gin-gonic/gin"
	entity2 "mio/internal/pkg/model/entity"
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
func GetAuthAdmin(c *gin.Context) entity2.Admin {
	if admin, ok := c.Get("AuthAdmin"); ok {
		return admin.(entity2.Admin)
	}
	return entity2.Admin{}
}
func GetAuthUser(c *gin.Context) entity2.User {
	if user, ok := c.Get("AuthUser"); ok {
		return user.(entity2.User)
	}
	return entity2.User{}
}
