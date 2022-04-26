package auth

import (
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
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

func (srv WeappService) LoginByCode(code string) (*entity.User, string, error) {
	//调用java那边登陆接口
	result, err := util.DefaultHttp.OriginJson(javaLoginUrl, "POST", []byte(fmt.Sprintf(`{"code":"%s"}`, code)))
	if err != nil {
		return nil, "", err
	}

	//获取用户信息
	cookie := result.Response.Header.Get("Set-Cookie")
	whoAmiResult, err := util.DefaultHttp.OriginGet(javaWhoAmi, util.HttpWithHeader("Cookie", cookie))
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
