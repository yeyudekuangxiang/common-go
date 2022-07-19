package auth

import (
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/httputil"
	"time"
)

const (
	javaLoginUrl = "https://dev.mp2c.miotech.com/api/mp2c/auth/login"
	javaWhoAmi   = "https://dev.mp2c.miotech.com/api/mp2c/user/whoami"
)

var userDealPool, _ = ants.NewPool(100)
var DefaultWeappService = WeappService{}

type WeappService struct {
	client *weapp.Client
}

func (srv WeappService) LoginByCode(code string, invitedBy string, partnershipWith entity.PartnershipType) (*entity.User, string, error) {
	//调用java那边登陆接口
	result, err := httputil.OriginJson(javaLoginUrl, "POST", []byte(fmt.Sprintf(`{"code":"%s"}`, code)))
	if err != nil {
		return nil, "", err
	}

	//获取用户信息
	cookie := result.Response.Header.Get("Set-Cookie")

	log.Println("cookie", cookie, invitedBy, partnershipWith)
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

	isNewUser := false
	if user.ID == 0 {
		isNewUser = true
		user, err = service.DefaultUserService.CreateUser(service.CreateUserParam{
			OpenId:      whoAmiResp.Data.Openid,
			AvatarUrl:   "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/images/user/default.png",
			Nickname:    "绿喵用户" + util.RandomStr(6, util.RandomStrNumber, util.RandomStrLower),
			PhoneNumber: "",
			Source:      entity.UserSourceMio,
			UnionId:     session.WxUnionId,
		})
		if err != nil {
			return nil, "", err
		}
		//return user, cookie, nil
	} else if user.GUID == "" && session.WxUnionId != "" { //更新用户unionid
		service.DefaultUserService.UpdateUserUnionId(user.ID, session.WxUnionId)

	}

	if isNewUser {
		err := userDealPool.Submit(func() {
			srv.AfterCreateUser(user, invitedBy, partnershipWith)
		})
		if err != nil {
			app.Logger.Errorf("提交新用户处理事件失败 %+v %s %s", user, invitedBy, partnershipWith)
		}
	}

	return user, cookie, nil
}

func (srv WeappService) AfterCreateUser(user *entity.User, invitedBy string, partnershipType entity.PartnershipType) {
	app.Logger.Infof("新用户创建后事件 %+v %s %s", user, invitedBy, partnershipType)
	_, err := service.DefaultStepService.FindOrCreateStep(user.ID)
	if err != nil {
		app.Logger.Error(user, invitedBy, err)
	}

	if invitedBy != "" {
		err := service.DefaultInviteService.Invite(user.OpenId, invitedBy)
		if err != nil {
			app.Logger.Errorf("添加邀请关系失败 %+v %s %s %v", user, invitedBy, partnershipType, err)
		} else {
			app.Logger.Errorf("添加邀请关系成功 %+v %s %s ", user, invitedBy, partnershipType)
		}
	}

	if partnershipType != "" {
		list, err := service.DefaultPartnershipRedemptionService.ProcessPromotionInformation(user.OpenId, partnershipType, entity.PartnershipPromotionTriggerREGISTER)
		if err != nil {
			app.Logger.Errorf("添加第三方活动信息失败 %+v %s %s %v", user, invitedBy, partnershipType, err)
		} else {
			app.Logger.Errorf("添加第三方活动信息成功 %+v %s %s %+v", user, invitedBy, partnershipType, list)
		}
	}
}
