package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/service"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	var form GetUserForm
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	user, err := service.DefaultUserService.GetUserById(form.Id)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}
