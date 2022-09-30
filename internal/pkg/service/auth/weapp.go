package auth

import (
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	activityService "mio/internal/pkg/service/activity"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/httputil"
	"time"
)

var userDealPool, _ = ants.NewPool(100)
var DefaultWeappService = WeappService{}

type WeappService struct {
	client *weapp.Client
}

func (srv WeappService) LoginByCode(code string, invitedBy string, partnershipWith entity.PartnershipType, cid int64, thirdId string) (*entity.User, string, bool, error) {
	//调用java那边登陆接口
	result, err := httputil.OriginJson(config.Config.Java.JavaLoginUrl, "POST", []byte(fmt.Sprintf(`{"code":"%s"}`, code)))
	if err != nil {
		return nil, "", false, err
	}

	//获取用户信息
	cookie := result.Response.Header.Get("Set-Cookie")

	app.Logger.Debug("cookie", cookie, invitedBy, partnershipWith)
	whoAmiResult, err := httputil.OriginGet(config.Config.Java.JavaWhoAmi, httputil.HttpWithHeader("Cookie", cookie))
	if err != nil {
		return nil, "", false, err
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
		return nil, "", false, errors.WithStack(err)
	}

	if whoAmiResp.Code != "success" {
		return nil, "", false, errors.New(whoAmiResp.Message)
	}

	user, err := service.DefaultUserService.GetUserByOpenId(whoAmiResp.Data.Openid)
	if err != nil {
		return nil, "", false, err
	}
	session, _ := service.DefaultSessionService.FindSessionByOpenId(whoAmiResp.Data.Openid)

	isNewUser := false
	if user.ID == 0 {
		isNewUser = true
		user, err = service.DefaultUserService.CreateUser(service.CreateUserParam{
			OpenId:      whoAmiResp.Data.Openid,
			AvatarUrl:   "https://resources.miotech.com/static/mp2c/images/user/default.png",
			Nickname:    "绿喵用户" + util.RandomStr(6, util.RandomStrNumber, util.RandomStrLower),
			PhoneNumber: "",
			Source:      entity.UserSourceMio,
			UnionId:     session.WxUnionId,
			ChannelId:   cid,
		})
		if err != nil {
			return nil, "", false, err
		}

		if thirdId != "" && cid == 1057 {
			activityService.NewZyhService(context.NewMioContext()).Create(srv_types.GetZyhGetInfoByDTO{
				Openid: whoAmiResp.Data.Openid,
				VolId:  thirdId,
			})
		}
	} else if user.GUID == "" && session.WxUnionId != "" { //更新用户unionid
		service.DefaultUserService.UpdateUserUnionId(user.ID, session.WxUnionId)
	}
	if isNewUser {
		err := userDealPool.Submit(func() {
			srv.AfterCreateUser(user, invitedBy, partnershipWith)
			//注册领积分
			srv.ReceivePoint(user)
		})
		if err != nil {
			app.Logger.Errorf("提交新用户处理事件失败 %+v %s %s", user, invitedBy, partnershipWith)
		}
	}

	return user, cookie, isNewUser, nil
}

func (srv WeappService) ToZhuGe(openId string, attr map[string]interface{}, eventName string) {
	go track.DefaultZhuGeService().Track(eventName, openId, attr)
}

func (srv WeappService) AfterCreateUser(user *entity.User, invitedBy string, partnershipType entity.PartnershipType) {
	app.Logger.Infof("新用户创建后事件 %+v %s %s", user, invitedBy, partnershipType)
	_, err := service.DefaultStepService.FindOrCreateStep(user.OpenId)
	if err != nil {
		app.Logger.Error(user, invitedBy, err)
	}

	if invitedBy != "" {
		err := service.DefaultInviteService.Invite(user.OpenId, invitedBy)
		if err != nil {
			app.Logger.Error(user, invitedBy, err)
		}
		println(err)
		//进入好友关系表
		_, errFriend := service.DefaultUserFriendService.Create(user, invitedBy)
		if errFriend != nil {
			app.Logger.Error(user, invitedBy, errFriend)
		}
		zhuGeAttr := make(map[string]interface{}, 0)
		zhuGeAttr["邀请人"] = invitedBy
		zhuGeAttr["用户"] = user.OpenId
		track.DefaultZhuGeService().Track(config.ZhuGeEventName.UserInvitedBy, user.OpenId, zhuGeAttr)
	}

	if partnershipType != "" {
		_, err := service.DefaultPartnershipRedemptionService.ProcessPromotionInformation(user.OpenId, partnershipType, entity.PartnershipPromotionTriggerREGISTER)
		if err != nil {
			app.Logger.Errorf("添加第三方活动信息失败 %+v %s %s %v", user, invitedBy, partnershipType, err)
		}
	}
}

func (srv WeappService) ReceivePoint(user *entity.User) {
	//获取渠道信息
	chInfo := service.DefaultUserChannelService.GetChannelByCid(user.ChannelId)
	//判断该渠道是否可领取积分
	if transactionType, ok := entity.PlatformMethodMap[chInfo.Code]; ok {
		point := entity.PointCollectValueMap[transactionType]
		_, err := service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:      user.OpenId,
			Type:        transactionType,
			BizId:       util.UUID(),
			ChangePoint: int64(point),
			AdminId:     0,
			Note:        transactionType.Text(),
		})
		if err != nil {
			app.Logger.Errorf("注册领取积分失败 用户ID: %d; 渠道id:%d; 失败原因:%s\n", user.ID, user.ChannelId, err.Error())
		}
	}
}
