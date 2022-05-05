package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	var form GetUserForm
	if err := apiutil.BindForm(c, &form); err != nil {
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
