package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultInviteController = InviteController{}

type InviteController struct {
}

func (InviteController) GetShareQrCode(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	info, err := service.DefaultInviteService.GetInviteQrCode(user.OpenId)
	return gin.H{
		"qrInfo": info,
	}, err
}
func (InviteController) GetInviteList(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	list, err := service.DefaultInviteService.GetInviteList(user.OpenId)
	return gin.H{
		"list": list,
	}, err
}
