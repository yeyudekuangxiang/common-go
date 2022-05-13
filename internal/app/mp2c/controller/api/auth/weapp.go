package auth

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/auth"
	"mio/internal/pkg/util/apiutil"
	"strings"
)

var DefaultWeappController = WeappController{}

type WeappController struct {
}

func (ctr WeappController) LoginByCode(ctx *gin.Context) (gin.H, error) {
	form := WeappAuthForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	partnershipWith := entity.PartnershipType(strings.ToUpper(form.PartnershipWith))
	user, cookie, err := auth.DefaultWeappService.LoginByCode(form.Code, form.InvitedBy, partnershipWith)
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
