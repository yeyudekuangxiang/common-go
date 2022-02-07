package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/model/request"
	"mio/service"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	var form request.ApiUserId
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	user, err := service.DefaultUserService.GetUserById(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}
