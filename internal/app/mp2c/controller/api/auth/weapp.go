package auth

import (
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/auth"
	"mio/internal/pkg/util/apiutil"
	"strings"
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
	user, cookie, err := auth.DefaultWeappService.LoginByCode(form.Code, form.InvitedBy, partnershipWith, form.Cid)
	if err != nil {
		zhuGeAttr["失败原因"] = err.Error()
		auth.DefaultWeappService.ToZhuGe("无openid", zhuGeAttr, config.ZhuGeEventName.UserLoginErr)
		return nil, err
	}
	token, err := service.DefaultUserService.CreateUserToken(user.ID)
	if err != nil {
		zhuGeAttr["失败原因"] = err.Error()
		auth.DefaultWeappService.ToZhuGe(user.OpenId, zhuGeAttr, config.ZhuGeEventName.UserLoginErr)
		return nil, err
	}
	auth.DefaultWeappService.ToZhuGe(user.OpenId, zhuGeAttr, config.ZhuGeEventName.UserLoginSuc)
	zhuGeUserIdentifyAttr := make(map[string]interface{}, 0)
	zhuGeUserIdentifyAttr["id"] = user.ID
	zhuGeUserIdentifyAttr["openid"] = user.OpenId
	zhuGeUserIdentifyAttr["性别"] = user.Gender
	zhuGeUserIdentifyAttr["注册来源"] = user.Source
	zhuGeUserIdentifyAttr["注册时间"] = user.Time.Format("2006/01/02")
	zhuGeUserIdentifyAttr["注册定位城市"] = user.CityCode
	zhuGeUserIdentifyAttr["用户渠道分类"] = user.ChannelId
	zhuGeUserIdentifyAttr["子渠道"] = user.ChannelId
	zhuGeUserIdentifyAttr["ip"] = user.Ip
	auth.DefaultWeappService.ToZhuGe(user.OpenId, zhuGeAttr, config.ZhuGeEventName.UserIdentify)

	return gin.H{
		"token":  token,
		"cookie": cookie,
	}, nil
}
