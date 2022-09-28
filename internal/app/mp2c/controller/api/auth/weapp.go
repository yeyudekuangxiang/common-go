package auth

import (
	"mio/config"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/auth"
	"mio/internal/pkg/util/apiutil"
	"strings"

	"github.com/gin-gonic/gin"
)

var DefaultWeappController = WeappController{}

type WeappController struct {
}

type TrackLoginZhuGe struct {
	OpenId string
	Event  string
}

func (ctr WeappController) LoginByCode(ctx *gin.Context) (gin.H, error) {
	form := WeappAuthForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	zhuGeAttr := make(map[string]interface{}, 0)
	partnershipWith := entity.PartnershipType(strings.ToUpper(form.PartnershipWith))
	user, cookie, isNewUser, err := auth.DefaultWeappService.LoginByCode(form.Code, form.InvitedBy, partnershipWith, form.Cid)
	if err != nil {
		zhuGeAttr["失败原因"] = err.Error()
		auth.DefaultWeappService.ToZhuGe("无openid", zhuGeAttr, config.ZhuGeEventName.UserLoginErr)
		return nil, err
	}

	go service.DefaultUserService.SendUserIdentifyToZhuGe(user.OpenId) //用户基本信息诸葛打点

	token, err := service.DefaultUserService.CreateUserToken(user.ID)
	if err != nil {
		zhuGeAttr["失败原因"] = err.Error()
		auth.DefaultWeappService.ToZhuGe(user.OpenId, zhuGeAttr, config.ZhuGeEventName.UserLoginErr)
		return nil, err
	}
	auth.DefaultWeappService.ToZhuGe(user.OpenId, zhuGeAttr, config.ZhuGeEventName.UserLoginSuc)

	return gin.H{
		"token":     token,
		"cookie":    cookie,
		"isNewUser": isNewUser,
	}, nil
}
