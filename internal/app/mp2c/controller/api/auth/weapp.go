package auth

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/auth"
	"mio/internal/pkg/service/srv_types"
	utilPkg "mio/internal/pkg/util"
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

	partnershipWith := entity.PartnershipType(strings.ToUpper(form.PartnershipWith))
	user, cookie, err := auth.DefaultWeappService.LoginByCode(form.Code, form.InvitedBy, partnershipWith, form.Cid)
	if err != nil {
		trackUser(TrackLoginZhuGe{OpenId: user.OpenId, Event: "登录失败"}, err.Error()) //失败了，诸葛打点
		return nil, err
	}
	token, err := service.DefaultUserService.CreateUserToken(user.ID)
	if err != nil {
		trackUser(TrackLoginZhuGe{OpenId: user.OpenId, Event: "登录失败"}, err.Error()) //失败了，诸葛打点
		return nil, err
	}
	trackUser(TrackLoginZhuGe{OpenId: user.OpenId, Event: "登录成功"}, "") //成功了，诸葛打点
	return gin.H{
		"token":  token,
		"cookie": cookie,
	}, nil
}

//用户打点
func trackUser(dto TrackLoginZhuGe, failMessage string) {
	service.DefaultZhuGeService().TrackLogin(srv_types.TrackLoginZhuGe{
		OpenId:      dto.OpenId,
		IsFail:      utilPkg.Ternary(failMessage == "", false, true).Bool(),
		FailMessage: failMessage,
		Event:       dto.Event,
	})
}
