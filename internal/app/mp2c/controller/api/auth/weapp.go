package auth

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/auth"
	"mio/internal/pkg/util"
)

var DefaultWeappController = WeappController{}

type WeappController struct {
}

func (ctr WeappController) LoginByCode(ctx *gin.Context) (gin.H, error) {
	form := WeappAuthForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user, cookie, err := auth.DefaultWeappService.LoginByCode(form.Code)
	if err != nil {
		return nil, err
	}
	token, err := service.DefaultUserService.CreateUserToken(user.ID)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"token":  token,
		"cookie": cookie,
	}, nil
}
