package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultInviteController = InviteController{}

type InviteController struct {
}

func (InviteController) GetShareQrCode(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	//page := "/pages/home/index?invitedBy=" + user.OpenId + "&cid=2"
	page := "/pages/invitation/index?invitedBy=" + user.OpenId + "&cid=2"
	info, err := service.NewQRCodeService().GetLimitedQRCode(entity.QrCodeSceneInvite, page, 100, user.OpenId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"qrcode": service.DefaultOssService.FullUrl(info.ImagePath),
	}, nil
}
func (InviteController) GetInviteList(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	list, err := service.DefaultInviteService.GetInviteList(user.OpenId)

	return gin.H{
		"list": list,
	}, err
}
