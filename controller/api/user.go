package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/service"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetNewUser(c *gin.Context) (gin.H, error) {
	user, err := service.DefaultUserService.GetUserById(1)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	user := util.GetAuthUser(c)
	return gin.H{
		"user": user,
	}, nil
}
