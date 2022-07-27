package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCheckinController = CheckinController{}

type CheckinController struct {
}

func (CheckinController) Checkin(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	day, err := service.DefaultCheckinService.Checkin(user.OpenId)
	return gin.H{
		"checkedNumber": day,
	}, err
}

func (CheckinController) GetCheckinInfo(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	info, err := service.DefaultCheckinService.GetCheckInfo(user.OpenId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"checkinInfo": info,
	}, nil
}
