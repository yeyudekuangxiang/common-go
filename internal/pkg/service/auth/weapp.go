package auth

import (
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/httputil"
	"time"
)

const (
	javaLoginUrl = "https://miniprogram.api.miotech.com/api/mp2c/auth/login"
	javaWhoAmi   = "https://miniprogram.api.miotech.com/api/mp2c/user/whoami"
)

var DefaultWeappService = WeappService{}

type WeappService struct {
	client *weapp.Client
}

func (srv WeappService) LoginByCode(code string, invitedBy string) (*entity.User, string, error) {
	//调用java那边登陆接口
	result, err := httputil.OriginJson(javaLoginUrl, "POST", []byte(fmt.Sprintf(`{"code":"%s"}`, code)))
	if err != nil {
		return nil, "", err
	}

	//获取用户信息
	cookie := result.Response.Header.Get("Set-Cookie")
	whoAmiResult, err := httputil.OriginGet(javaWhoAmi, httputil.HttpWithHeader("Cookie", cookie))
	if err != nil {
		return nil, "", err
	}
	whoAmiResp := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Openid         string `json:"openid"`
			Registered     bool   `json:"registered"`
			IsAdmin        bool   `json:"isAdmin"`
			HasPhoneNumber bool   `json:"hasPhoneNumber"`
		} `json:"data"`
		ResponseAt time.Time `json:"responseAt"`
	}{}

	err = json.Unmarshal(whoAmiResult.Body, &whoAmiResp)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	if whoAmiResp.Code != "success" {
		return nil, "", errors.New(whoAmiResp.Message)
	}

	user, err := service.DefaultUserService.GetUserByOpenId(whoAmiResp.Data.Openid)
	if err != nil {
		return nil, "", err
	}
	session, _ := service.DefaultSessionService.FindSessionByOpenId(whoAmiResp.Data.Openid)

	if user.ID == 0 {
		user, err := service.DefaultUserService.CreateUser(service.CreateUserParam{
			OpenId:      whoAmiResp.Data.Openid,
			AvatarUrl:   "",
			Nickname:    "",
			PhoneNumber: "",
			Source:      entity.UserSourceMio,
			UnionId:     session.WxUnionId,
		})
		if err != nil {
			return nil, "", err
		}
		return user, cookie, nil
	} else if user.GUID == "" && session.WxUnionId != "" { //更新用户unionid
		service.DefaultUserService.UpdateUserUnionId(user.ID, session.WxUnionId)
	}

	return user, cookie, nil
}

func (srv WeappService) AfterCreateUser(user *entity.User, invitedBy string, partnershipType entity.PartnershipType) {
	_, err := service.DefaultStepService.FindOrCreateStep(user.ID)
	if err != nil {
		app.Logger.Error(user, invitedBy, err)
	}

	if invitedBy != "" {
		_, isNew, err := service.DefaultInviteService.AddInvite(user.OpenId, invitedBy)
		if err != nil {
			app.Logger.Error(user, invitedBy, err)
		} else if isNew {
			//发放积分奖励
			pointService := service.NewPointService(context.NewMioContext())
			_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
				OpenId:       invitedBy,
				Type:         entity.POINT_INVITE,
				ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_INVITE]),
				BizId:        util.UUID(),
				AdditionInfo: fmt.Sprintf("invite %s", user.OpenId),
			})
			if err != nil {
				app.Logger.Error("发放邀请积分失败", err)
			}
		}
	}

	if partnershipType != "" {

	}
}
