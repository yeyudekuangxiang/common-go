package service

import (
	"context"
	"fmt"
	"github.com/medivhzhan/weapp/v3/phonenumber"
	"math/rand"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	contextMix "mio/internal/pkg/core/context"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	util2 "mio/internal/pkg/util"
	"mio/internal/pkg/util/message"
	"mio/pkg/baidu"
	"mio/pkg/errno"
	"mio/pkg/wxapp"
	"strconv"
	"strings"
	"time"
)

var DefaultUserService = NewUserService(
	repository.DefaultUserRepository,
	repository.InviteRepository{},
	repository.NewCityRepository(contextMix.NewMioContext()),
	repository.NewUserExtendInfoRepository(contextMix.NewMioContext()))

func NewUserService(r repository.UserRepository,
	rInvite repository.InviteRepository, rCity repository.CityRepository, rUserExtend repository.UserExtentInfoRepository) UserService {
	return UserService{
		r:           r,
		rInvite:     rInvite,
		rCity:       rCity,
		rUserExtend: rUserExtend,
	}
}

type UserService struct {
	r           repository.UserRepository
	rInvite     repository.InviteRepository
	rCity       repository.CityRepository
	rUserExtend repository.UserExtentInfoRepository
}

// GetUser 查询用户信息
func (u UserService) GetUser(by repository.GetUserBy) (*entity.User, bool, error) {
	return u.r.GetUser(by)
}

// GetUserExtend 查询用户信息

func (u UserService) GetUserExtend(by repository.GetUserExtendBy) (*entity.UserExtendInfo, bool, error) {
	return u.rUserExtend.GetUserExtend(by)
}

func (u UserService) CreateUserExtend(param CreateUserExtendParam) (*entity.UserExtendInfo, error) {
	user, exist, err := u.rUserExtend.GetUserExtend(repository.GetUserExtendBy{
		OpenId: param.OpenId,
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		userExtend := &entity.UserExtendInfo{}
		if err := util2.MapTo(param, &userExtend); err != nil {
			return nil, err
		}
		city, errCity := baidu.IpToCity(param.Ip)
		if errCity != nil {
			app.Logger.Info("BindPhoneByCode ip地址查询失败", err.Error())
		}

		userExtend.CreatedAt = time.Now()
		userExtend.Adcode = city.Content.AddressDetail.Adcode
		userExtend.CityCode = city.Content.AddressDetail.CityCode
		userExtend.Province = city.Content.AddressDetail.Province
		userExtend.City = city.Content.AddressDetail.City
		userExtend.District = city.Content.AddressDetail.District
		userExtend.Street = city.Content.AddressDetail.Street
		userExtend.StreetNumber = city.Content.AddressDetail.StreetNumber

		ret := u.rUserExtend.Create(userExtend)
		return userExtend, ret
	} else {
		if param.Ip != "" {
			user.Ip = param.Ip
		}
		ret := u.rUserExtend.Save(user)
		return user, ret
	}
}

// GetUserByID 根据用户id获取用户信息
func (u UserService) GetUserByID(id int64) (*entity.User, bool, error) {
	return u.r.GetUserByID(id)
}

// GetUserById 请使用 GetUserByID 代替此方法
func (u UserService) GetUserById(id int64) (*entity.User, error) {
	user := u.r.GetUserById(id)
	return &user, nil
}
func (u UserService) GetUserByOpenId(openId string) (*entity.User, error) {
	if openId == "" {
		return &entity.User{}, nil
	}
	user := u.r.GetUserBy(repository.GetUserBy{
		OpenId: openId,
	})
	return &user, nil
}
func (u UserService) GetUserByToken(token string) (*entity.User, error) {
	var authUser auth.User
	err := util2.ParseToken(token, &authUser)
	if err != nil {
		return nil, err
	}
	user := u.r.GetUserById(authUser.ID)
	return &user, nil
}
func (u UserService) CreateUserToken(id int64) (string, error) {
	if id == 0 {
		return "", errno.ErrUserNotFound.WithCaller()
	}

	user := u.r.GetUserById(id)
	if user.ID == 0 {
		return "", errno.ErrUserNotFound.WithCaller()
	}

	return util2.CreateToken(auth.User{
		ID:        user.ID,
		Mobile:    user.PhoneNumber,
		CreatedAt: model.Time{Time: time.Now()},
	})
}

//SendUserIdentifyToZhuGe 用户属性上报到诸葛
func (u UserService) SendUserIdentifyToZhuGe(openid string) {
	return
	if openid == "" {
		return
	}
	user, exit, _ := u.r.GetUserIdentifyInfo(openid)
	if !exit {
		return //不存在用户信息，返回
	}
	zhuGeIdentifyAttr := make(map[string]interface{}, 0)
	zhuGeIdentifyAttr["用户渠道分类"] = user.ChannelTypeName
	zhuGeIdentifyAttr["子渠道"] = user.ChannelName
	track.DefaultZhuGeService().Track(config.ZhuGeEventName.UserIdentify, openid, zhuGeIdentifyAttr)
}

//func (u UserService) SendUserRegisterCoupon(user entity.User) {
//	orderNo := "jhx" + strconv.FormatInt(time.Now().Unix(), 10)
//	startTime, _ := time.Parse("2006-01-02", "2022-09-29")
//	endTime, _ := time.Parse("2006-01-02", "2022-10-31")
//	for i := 0; i < 2; i++ {
//		err := jhx.NewJhxService(mioctx.NewMioContext()).TicketCreate(orderNo+strconv.Itoa(i), 123, startTime, endTime, user)
//		if err != nil {
//			return
//		}
//	}
//	return
//}

func (u UserService) CreateUser(param CreateUserParam) (*entity.User, error) {
	user := u.r.GetUserBy(repository.GetUserBy{
		OpenId: param.OpenId,
		Source: param.Source,
	})
	if user.ID != 0 {
		return &user, nil
	}

	user = entity.User{}
	if err := util2.MapTo(param, &user); err != nil {
		return nil, err
	}
	user.Time = model.NewTime()
	user.Auth = 1
	user.Position = "ordinary"
	user.Partners = 2
	user.Status = 1

	ch := DefaultUserChannelService.GetChannelByCid(param.ChannelId) //获取渠道id
	user.ChannelId = ch.Cid
	ret := repository.DefaultUserRepository.Save(&user)
	/*
		retCity, cityErr := u.rCity.GetByCityCode(repotypes.GetCityByCode{CityCode: "140900"})
	*/
	/*	userExtend, exist, _ := u.rUserExtend.GetUserExtend(repository.GetUserExtendBy{
			OpenId: user.OpenId,
		})
	*/
	//上报到诸葛
	zhuGeAttr := make(map[string]interface{}, 0)
	zhuGeAttr["来源"] = param.Source
	zhuGeAttr["渠道"] = ch.Name
	zhuGeAttr["城市code"] = user.CityCode
	zhuGeAttr["openid"] = user.OpenId
	zhuGeAttr["ip"] = user.Ip
	zhuGeAttr["省"] = param.Province
	zhuGeAttr["市"] = param.City

	/*if cityErr == nil {
		zhuGeAttr["城市名"] = retCity.Name
	}*/
	if ret == nil {
		zhuGeAttr["是否成功"] = "成功"
	} else {
		zhuGeAttr["是否成功"] = "失败"
		zhuGeAttr["失败原因"] = ret.Error()
	}

	track.DefaultZhuGeService().Track(config.ZhuGeEventName.NewUserAdd, param.OpenId, zhuGeAttr)
	return &user, ret
}
func (u UserService) UpdateUserUnionId(id int64, unionid string) {
	if unionid == "" || id == 0 {
		return
	}
	user := u.r.GetUserById(id)
	if user.ID == 0 {
		return
	}
	app.DB.Model(entity.User{}).Where("id = ?", id).UpdateColumn("unionid = ?", unionid)
}

// GetUserBy 请使用 GetUser 代替此方法
func (u UserService) GetUserBy(by repository.GetUserBy) (*entity.User, error) {
	user := repository.DefaultUserRepository.GetUserBy(by)
	return &user, nil
}
func (u UserService) FindOrCreateByMobile(mobile string, cid int64) (*entity.User, error) {
	if mobile == "" {
		return nil, errno.ErrCommon.WithMessage("手机号不能为空")
	}
	user := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: mobile,
		Source: entity.UserSourceMobile,
	})
	if user.ID > 0 {
		return &user, nil
	}
	ch := DefaultUserChannelService.GetChannelByCid(cid) //获取渠道来源
	return u.CreateUser(CreateUserParam{
		OpenId:      mobile,
		Nickname:    "手机用户" + mobile[len(mobile)-4:],
		PhoneNumber: mobile,
		Source:      entity.UserSourceMobile,
		UnionId:     mobile,
		ChannelId:   ch.Cid,
	})
}

// BindMobileByYZM 绑定手机号
func (u UserService) BindMobileByYZM(userId int64, mobile string) error {
	if mobile == "" {
		return errno.ErrCommon.WithMessage("手机号不能为空")
	}
	userBy := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{Mobile: mobile})
	if userBy.ID != 0 {
		return errno.ErrCommon.WithMessage("该号码已被绑定，请更换号码重新绑定")
	}
	user := repository.DefaultUserRepository.GetUserById(userId)
	if user.ID == 0 {
		return errno.ErrCommon.WithMessage("未查到用户信息")
	}
	if user.PhoneNumber != "" {
		return errno.ErrCommon.WithMessage("您已绑定号码，请勿重复操作")
	}
	user.PhoneNumber = mobile
	return repository.DefaultUserRepository.Save(&user)
}

// FindUserBySource 根据用户id 获取指定平台的用户
func (u UserService) FindUserBySource(source entity.UserSource, userId int64) (*entity.User, error) {
	if userId == 0 {
		return &entity.User{}, nil
	}

	user := repository.DefaultUserRepository.GetUserById(userId)

	if user.ID == 0 || user.PhoneNumber == "" {
		return &entity.User{}, nil
	}

	sourceUer := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: user.PhoneNumber,
		Source: source,
	})

	return &sourceUer, nil
}
func (u UserService) GetYZM(mobile string) (string, error) {
	code := ""
	for i := 0; i < 4; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}
	//加入缓存
	cmd := app.Redis.Set(context.Background(), config.RedisKey.YZM+mobile, code, time.Second*30*60)
	fmt.Println(cmd.String())
	//发送短信
	message.SendYZM(mobile, code)

	return code, nil
}
func (u UserService) GetYZM2B(mobile string) (string, error) {
	code := ""
	for i := 0; i < 4; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}
	//加入缓存
	cmd := app.Redis.Set(context.Background(), config.RedisKey.YZM2B+mobile, code, time.Second*10*60)
	fmt.Println(cmd.String())
	//发送短信
	_, err := message.SendYZMSms2B(mobile, code)
	if err != nil {
		return "", err
	}
	return code, nil
}
func (u UserService) CheckYZM(mobile string, code string) bool {
	//取出缓存
	codeCmd := app.Redis.Get(context.Background(), config.RedisKey.YZM+mobile)
	fmt.Println(codeCmd.String())
	if codeCmd.Val() == code {
		fmt.Println("验证码验证通过")
		return true
	}
	return false
}

//企业版验证验证码

func (u UserService) CheckYZM2B(mobile string, code string) bool {
	//取出缓存
	codeCmd := app.Redis.Get(context.Background(), config.RedisKey.YZM2B+mobile)
	if codeCmd.Val() == code {
		fmt.Println("验证码验证通过")
		return true
	}
	return false
}

func (u UserService) BindPhoneByCode(userId int64, code string, cip string, invitedBy string) error {
	userInfo := u.r.GetUserById(userId)

	if userInfo.ID == 0 {
		app.Logger.Errorf("BindPhoneByCode Error:%s userId:%d", "未查到用户信息", userId)
		return errno.ErrCommon.WithMessage("未查到用户信息")
	}

	var phoneResult *phonenumber.GetPhoneNumberResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		phoneResult, err = app.Weapp.NewPhonenumber().GetPhoneNumber(&phonenumber.GetPhoneNumberRequest{
			Code: code,
		})
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(phoneResult.ErrCode)
	}, 1)

	if err != nil {
		app.Logger.Errorf("BindPhoneByCode error:%s", err.Error())
		return errno.ErrCommon
	}

	if phoneResult.ErrCode != 0 {
		app.Logger.Errorf("BindPhoneByCode Response error:%s", phoneResult.ErrMSG)
		return errno.ErrBindMobile.WithErrMessage(fmt.Sprintf("%d %s", phoneResult.ErrCode, phoneResult.ErrMSG))
	}

	userInfo.PhoneNumber = phoneResult.Data.PhoneNumber

	isBlack := app.Redis.SIsMember(context.Background(), config.RedisKey.BlackList, phoneResult.Data.PhoneNumber)
	if isBlack.Val() {
		userInfo.Risk = 4
	} else {
		//检测用户风险等级
		userRiskRankParam := wxapp.UserRiskRankParam{
			AppId:    config.Config.Weapp.AppId,
			OpenId:   userInfo.OpenId,
			Scene:    0,
			ClientIp: cip,
			MobileNo: userInfo.PhoneNumber,
		}
		var rest *wxapp.UserRiskRankResponse
		err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
			rest, err = app.Weapp.GetUserRiskRank(userRiskRankParam)
			if err != nil {
				return false, err
			}
			return app.Weapp.IsExpireAccessToken(rest.ErrCode)
		}, 1)

		if err != nil {
			app.Logger.Info("BindPhoneByCode 风险等级查询查询出错", err.Error())
		}

		userInfo.Risk = rest.RiskRank
	}

	//获取用户地址  todo 加入队列
	city, err := baidu.IpToCity(cip)
	if err != nil {
		app.Logger.Errorf("BindPhoneByCode ip地址查询失败 %s", err.Error())
	} else {
		userInfo.CityCode = city.Content.AddressDetail.Adcode
		userInfo.Ip = cip
	}

	userByMobile, ok, _ := u.r.GetUser(repository.GetUserBy{Mobile: userInfo.PhoneNumber, Source: entity.UserSourceMio})
	specialUser := DefaultUserSpecialService.GetSpecialUserByPhone(userInfo.PhoneNumber)

	//检查重复绑定 特殊用户有已绑定的账号
	if ok && userByMobile.OpenId != userInfo.OpenId && specialUser.ID == 0 {
		app.Logger.Errorf("BindPhoneByCode err: bind user: %s; old user: %s, bind mobile:%s, binding mobile:%s, isSpecial:%d", userInfo.OpenId, userByMobile.OpenId, userByMobile.PhoneNumber, userInfo.PhoneNumber, specialUser.ID)
		return errno.ErrCommon.WithMessage("该号码已绑定")
	}

	//更新特殊用户的数据
	if ok && specialUser.ID != 0 && !u.checkOpenId(userByMobile.OpenId) && specialUser.Status == 0 {
		//更新topic userid
		err := community.DefaultTopicService.UpdateAuthor(userInfo.ID, userByMobile.ID)
		if err != nil {
			app.Logger.Info("special用户topic数据更新失败", err.Error())
		}
		userByMobile.PhoneNumber = ""
		err = u.r.Save(userByMobile)
		if err != nil {
			app.Logger.Info("special用户原账号更新失败", err.Error())
		}
		specialUser.Status = 1
		err = DefaultUserSpecialService.Save(&specialUser)
		if err != nil {
			app.Logger.Info("special用户状态更新失败", err.Error())
		}
		userInfo.Auth = 1
		userInfo.PositionIcon = specialUser.Icon
		userInfo.Partners = entity.Partner(specialUser.Partner)
		userInfo.Position = entity.UserPosition(specialUser.Identity)
	}

	//更新保存用户信息
	ret := u.r.Save(&userInfo)

	go u.SendUserIdentifyToZhuGe(userInfo.OpenId) //个人信息打点到诸葛

	if invitedBy != "" && userInfo.Risk > 2 {
		return ret
	}

	//有邀请，并且没有发放奖励，不是黑产用户，给用户发放奖励
	inviteInfo := u.rInvite.GetInviteNoReward(userInfo.OpenId)
	if inviteInfo.ID != 0 && userInfo.Risk <= 2 {
		//发放积分奖励
		_, err = NewPointService(mioctx.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       inviteInfo.InvitedByOpenId,
			Type:         entity.POINT_INVITE,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_INVITE]),
			AdditionInfo: fmt.Sprintf("invite %s", userInfo.OpenId),
			InviteId:     inviteInfo.ID,
		})
		if err != nil {
			app.Logger.Errorf("发放邀请积分失败:%s, 用户openId:%s", err.Error(), userInfo.OpenId)
		}

	}

	/*	//随申行，绑定关系
		_, err = app.RpcService.ActivityRpcSrv.UpdateActivityThirdUser(context.Background(), &activityclient.UpdateActivityThirdUserReq{
			ActivityId: 2,
			UserId:     userInfo.ID,
			Openid:     userInfo.OpenId,
			Phone:      userInfo.PhoneNumber,
		})
		if err != nil {
			app.Logger.Errorf("【绑定手机号】随申行绑定手机号失败:%s", err.Error())
		}*/
	return ret
}
func (u UserService) BindPhoneByIV(param BindPhoneByIVParam) error {
	userInfo := u.r.GetUserById(param.UserId)
	if userInfo.ID == 0 {
		return errno.ErrCommon.WithMessage("未查到用户信息")
	}

	session, err := DefaultSessionService.FindSessionByOpenId(userInfo.OpenId)

	if err != nil {
		return err
	}
	if session.ID == 0 {
		return errno.ErrAuth.WithCaller()
	}

	decryptResult, err := app.Weapp.DecryptMobile(session.SessionKey, param.EncryptedData, param.IV)

	if err != nil {
		return err
	}
	userInfo.PhoneNumber = decryptResult.PhoneNumber
	return u.r.Save(&userInfo)
}

type Summery struct {
	CurrentSteps        int     `json:"currentSteps"`
	RedeemedPointsToday int64   `json:"redeemedPointsToday"`
	BalanceOfPoints     int64   `json:"balanceOfPoints"`
	SavedCO2            float64 `json:"savedCO2"`
	PendingPoints       int     `json:"pendingPoints"`
	StepDiff            int     `json:"stepDiff"`
}

// UserSummary 获取用户步行相关的数据统计
func (u UserService) UserSummary(userId int64) (*Summery, error) {
	summery := Summery{}

	userInfo := u.r.GetUserById(userId)
	if userInfo.ID == 0 {
		return &summery, nil
	}

	lastStepHistory, err := DefaultStepHistoryService.FindStepHistory(FindStepHistoryBy{
		OpenId:  userInfo.OpenId,
		Day:     model.NewTime().StartOfDay(),
		OrderBy: entity.OrderByList{entity.OrderByStepHistoryCountDesc},
	})
	if err != nil {
		return nil, err
	}
	//今日步数
	summery.CurrentSteps = lastStepHistory.Count

	//获取用户今日已领取步行积分
	todayStepPoint, err := u.calculateStepPointsOfToday(userId)
	if err != nil {
		return nil, err
	}
	summery.RedeemedPointsToday = todayStepPoint

	pointService := NewPointService(mioctx.NewMioContext())
	point, err := pointService.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	summery.BalanceOfPoints = point.Balance

	summery.SavedCO2 = DefaultCarbonNeutralityService.calculateCO2ByStep(int64(lastStepHistory.Count))

	pendingPoints, _, err := DefaultStepService.ComputePendingPoint(userInfo.OpenId)
	if err != nil {
		return nil, err
	}
	summery.PendingPoints = int(pendingPoints)

	stepDiff, err := u.getStepDiffFromDates(userId, model.NewTime().StartOfDay(), model.Time{Time: model.NewTime().Add(-time.Hour * 24)}.StartOfDay())
	if err != nil {
		return nil, err
	}
	summery.StepDiff = stepDiff
	return &summery, nil
}

//获取用户今日已领取步行积分
func (u UserService) calculateStepPointsOfToday(userId int64) (int64, error) {
	userInfo := u.r.GetUserById(userId)
	if userInfo.ID == 0 {
		return 0, nil
	}

	t := model.NewTime()
	pointTranService := NewPointTransactionService(mioctx.NewMioContext())
	list := pointTranService.GetListBy(repository.GetPointTransactionListBy{
		OpenId:    userInfo.OpenId,
		StartTime: t.StartOfDay(),
		EndTime:   t.EndOfDay(),
		Type:      entity.POINT_STEP,
	})
	var total int64 = 0
	for _, item := range list {
		total += item.Value
	}
	return total, nil
}

// history 今天的步行历史 step 步行总历史
func (u UserService) computePendingHistoryStep(history entity.StepHistory, step entity.Step) int {
	// date check is moved outside
	lastCheckedSteps := 0
	stepUpperLimit := StepToScoreConvertRatio * StepScoreUpperLimit

	//如果最后一次领积分时间为0 或者 最后一次领取时间不等于今天的开始时间
	if step.LastCheckTime.IsZero() || !step.LastCheckTime.Equal(model.NewTime().StartOfDay().Time) {
		lastCheckedSteps = 0
	} else {
		lastCheckedSteps = step.LastCheckCount
		if lastCheckedSteps > stepUpperLimit {
			return 0
		}
	}

	currentStep := util2.Ternary(history.Count < stepUpperLimit, history.Count, stepUpperLimit).Int()
	result := currentStep - lastCheckedSteps + lastCheckedSteps%StepToScoreConvertRatio

	return util2.Ternary(result > 0, result, 0).Int()
}

//比昨天减少
func (u UserService) getStepDiffFromDates(userId int64, day1 model.Time, day2 model.Time) (int, error) {
	userinfo := u.r.GetUserById(userId)
	if userinfo.ID == 0 {
		return 0, nil
	}

	stepHistory1, err := DefaultStepHistoryService.FindStepHistory(FindStepHistoryBy{
		Day:    day1,
		OpenId: userinfo.OpenId,
	})
	if err != nil {
		return 0, err
	}
	stepHistory2, err := DefaultStepHistoryService.FindStepHistory(FindStepHistoryBy{
		Day:    day2,
		OpenId: userinfo.OpenId,
	})
	if err != nil {
		return 0, err
	}
	return stepHistory1.Count - stepHistory2.Count, nil
}
func (u UserService) GetUserListBy(by repository.GetUserListBy) ([]entity.User, error) {
	return u.r.GetUserListBy(by), nil
}
func (u UserService) UpdateUserInfo(param UpdateUserInfoParam) error {
	//图片审核
	user := u.r.GetUserById(param.UserId)

	if user.ID == 0 {
		return errno.ErrUserNotFound
	}

	if param.PhoneNumber != nil {
		if u.CheckMobileBound(entity.UserSourceMio, user.ID, *param.PhoneNumber) {
			return errno.ErrCommon.WithMessage("该手机号已被其他账号绑定")
		}
		user.PhoneNumber = *param.PhoneNumber
	}
	if param.Birthday != nil {
		user.Birthday = model.Date{Time: *param.Birthday}
	}
	if param.Gender != nil {
		user.Gender = *param.Gender
	}
	if param.Position != "" {
		user.Position = entity.UserPosition(param.Position)
	}
	if param.Partners != 0 {
		user.Partners = entity.Partner(param.Partners)
	}
	if param.Status != 0 {
		user.Status = param.Status
	}
	if param.Auth != 0 {
		user.Auth = param.Auth
	}
	if param.Nickname != "" {
		user.Nickname = param.Nickname
	}
	if param.Avatar != "" {
		user.AvatarUrl = param.Avatar
	}
	if param.Introduction != "" {
		user.Introduction = param.Introduction
	}
	return u.r.Save(&user)
}
func (u UserService) GetUserPageListBy(by repository.GetUserPageListBy) ([]entity.User, int64) {
	return u.r.GetUserPageListBy(by)
}

func (u UserService) UpdateUserRisk(param UpdateUserRiskParam) error {
	user := u.r.GetUserById(param.UserId)
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}
	user.Risk = param.Risk
	return u.r.Save(&user)
}

// CheckUserRisk 检测用户风险等级
func (u UserService) CheckUserRisk(param wxapp.UserRiskRankParam) (*wxapp.UserRiskRankResponse, error) {
	var rest *wxapp.UserRiskRankResponse
	err := app.Weapp.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		rest, err = app.Weapp.GetUserRiskRank(param)
		if err != nil {
			return false, err
		}
		return app.Weapp.IsExpireAccessToken(rest.ErrCode)
	}, 1)

	if err != nil {
		return nil, err
	}
	//创建记录
	err2 := DefaultUserRiskLogService.Create(&entity.UserRiskLog{
		OpenId:   param.OpenId,
		Scene:    param.Scene,
		MobileNo: param.MobileNo,
		ClientIp: param.ClientIp,
		ErrCode:  rest.ErrCode,
		ErrMsg:   rest.ErrMsg,
		UnoinId:  rest.UnoinId,
		RiskRank: rest.RiskRank,
	})
	if err2 != nil {
		app.Logger.Error("DefaultUserRiskLogService.Create 异常", rest)
	}
	return rest, err
}

// AccountInfo 用户账户信息
func (u UserService) AccountInfo(userId int64) (*UserAccountInfo, error) {
	point, err := NewPointService(mioctx.NewMioContext()).FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	go DefaultUserService.SendUserIdentifyToZhuGe(point.OpenId) //用户基本信息诸葛打点
	certCount, err := DefaultBadgeService.GetUserCertCountById(userId)
	if err != nil {
		return nil, err
	}
	carbonInfo, _ := NewCarbonTransactionService(mioctx.NewMioContext()).Info(api_types.GetCarbonTransactionInfoDto{UserId: userId})

	return &UserAccountInfo{
		Balance:     point.Balance,
		CertNum:     certCount,
		CarbonToday: carbonInfo.CarbonToday,
		CarbonAll:   carbonInfo.Carbon,
	}, nil
}

func (u UserService) CheckMobileBound(source entity.UserSource, id int64, mobile string) bool {
	user := u.r.GetUserBy(repository.GetUserBy{
		Source: source,
		Mobile: mobile,
	})

	if user.ID == 0 || user.ID == id {
		return false
	}
	return true
}

func (u UserService) ChangeUserState(param ChangeUserState) error {
	user := u.r.GetUserById(param.UserId)
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}
	if param.Status >= 0 {
		user.Status = param.Status
	}
	return u.r.Save(&user)
}

func (u UserService) ChangeUserPosition(param ChangeUserPosition) error {
	user := u.r.GetUserById(param.UserId)
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}
	if param.Position != "" {
		user.Position = entity.UserPosition(param.Position)
		if icon, ok := entity.IconMap[user.Position]; ok {
			user.PositionIcon = icon
		}
	}
	return u.r.Save(&user)
}

func (u UserService) ChangeUserPartner(param ChangeUserPartner) error {
	user := u.r.GetUserById(param.UserId)
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}
	if param.Partner >= 0 {
		user.Partners = entity.Partner(param.Partner)
		if icon, ok := entity.IconPartnerMap[user.Partners]; ok {
			user.PositionIcon = icon
		}
	}
	return u.r.Save(&user)
}

func (u UserService) checkOpenId(openId string) bool {
	return strings.HasPrefix(openId, "oy_")
}
