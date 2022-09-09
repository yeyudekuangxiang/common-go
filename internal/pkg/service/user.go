package service

import (
	"context"
	"fmt"
	"github.com/medivhzhan/weapp/v3/phonenumber"
	"github.com/pkg/errors"
	"math/rand"
	"mio/config"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	util2 "mio/internal/pkg/util"
	"mio/internal/pkg/util/message"
	"mio/pkg/baidu"
	"mio/pkg/errno"
	"mio/pkg/wxapp"
	"strconv"
	"time"
)

var DefaultUserService = NewUserService(repository.DefaultUserRepository, repository.InviteRepository{})

func NewUserService(r repository.UserRepository, rInvite repository.InviteRepository) UserService {
	return UserService{
		r:       r,
		rInvite: rInvite,
	}
}

type UserService struct {
	r       repository.UserRepository
	rInvite repository.InviteRepository
}

// GetUser 查询用户信息
func (u UserService) GetUser(by repository.GetUserBy) (*entity.User, bool, error) {
	return u.r.GetUser(by)
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
		return "", errno.ErrUserNotFound
	}

	user := u.r.GetUserById(id)
	if user.ID == 0 {
		return "", errno.ErrUserNotFound
	}

	return util2.CreateToken(auth.User{
		ID:        user.ID,
		Mobile:    user.PhoneNumber,
		CreatedAt: model.Time{Time: time.Now()},
	})
}

//SendUserIdentifyToZhuGe 用户属性上报到诸葛
func (u UserService) SendUserIdentifyToZhuGe(openid string) {
	user := u.r.GetUserIdentifyInfo(openid)
	zhuGeIdentifyAttr := make(map[string]interface{}, 0)
	zhuGeIdentifyAttr["openid"] = user.Openid
	zhuGeIdentifyAttr["注册来源"] = user.Source
	zhuGeIdentifyAttr["注册时间"] = user.Time.Format("2006/01/02")
	zhuGeIdentifyAttr["注册定位城市"] = user.CityName
	zhuGeIdentifyAttr["用户渠道分类"] = user.ChannelTypeName
	zhuGeIdentifyAttr["子渠道"] = user.ChannelName
	DefaultZhuGeService().Track(config.ZhuGeEventName.UserIdentify, openid, zhuGeIdentifyAttr)
}

func (u UserService) CreateUser(param CreateUserParam) (*entity.User, error) {
	user := u.r.GetUserBy(repository.GetUserBy{
		OpenId: param.OpenId,
		Source: param.Source,
	})
	if user.ID != 0 {
		return &user, nil
	}

	guid := ""
	if param.UnionId != "" {
		guid = u.r.GetGuid(param.UnionId)
		if guid == "" {
			guid = util2.UUID()
		}
	}

	user = entity.User{}
	if err := util2.MapTo(param, &user); err != nil {
		return nil, err
	}
	user.GUID = guid
	user.Time = model.NewTime()

	if param.UnionId != "" {
		app.DB.Where("unionid = ? and guid =''", param.UnionId).Update("guid", guid)
	}
	channelId := DefaultUserChannelService.GetChannelByCid(param.ChannelId) //获取渠道id
	user.ChannelId = channelId
	return &user, repository.DefaultUserRepository.Save(&user)
}
func (u UserService) UpdateUserUnionId(id int64, unionid string) {
	if unionid == "" || id == 0 {
		return
	}

	user := u.r.GetUserById(id)
	if user.ID == 0 {
		return
	}

	guid := u.r.GetGuid(unionid)
	if guid == "" {
		guid = util2.UUID()
	}

	user.UnionId = unionid
	user.GUID = guid
	err := u.r.Save(&user)
	if err != nil {
		app.Logger.Error("更新unionid失败", id, unionid, err)
	}
}

// GetUserBy 请使用 GetUser 代替此方法
func (u UserService) GetUserBy(by repository.GetUserBy) (*entity.User, error) {
	user := repository.DefaultUserRepository.GetUserBy(by)
	return &user, nil
}
func (u UserService) FindOrCreateByMobile(mobile string, cid int64) (*entity.User, error) {
	if mobile == "" {
		return nil, errors.New("手机号不能为空")
	}
	user := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{
		Mobile: mobile,
		Source: entity.UserSourceMobile,
	})
	if user.ID > 0 {
		return &user, nil
	}
	channelId := DefaultUserChannelService.GetChannelByCid(cid) //获取渠道来源
	return u.CreateUser(CreateUserParam{
		OpenId:      mobile,
		Nickname:    "手机用户" + mobile[len(mobile)-4:],
		PhoneNumber: mobile,
		Source:      entity.UserSourceMobile,
		UnionId:     mobile,
		ChannelId:   channelId,
	})
}

// BindMobileByYZM 绑定手机号
func (u UserService) BindMobileByYZM(userId int64, mobile string) error {
	if mobile == "" {
		return errors.New("手机号不能为空")
	}
	userBy := repository.DefaultUserRepository.GetUserBy(repository.GetUserBy{Mobile: mobile})
	if userBy.ID != 0 {
		return errors.New("该号码已被绑定，请更换号码重新绑定")
	}
	user := repository.DefaultUserRepository.GetUserById(userId)
	if user.ID == 0 {
		return errors.New("未查到用户信息")
	}
	if user.PhoneNumber != "" {
		return errors.New("您已绑定号码，请勿重复操作")
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
func (u UserService) BindPhoneByCode(userId int64, code string, cip string, invitedBy string) error {
	userInfo := u.r.GetUserById(userId)
	if userInfo.ID == 0 {
		return errors.New("未查到用户信息")
	}

	phoneResult, err := app.Weapp.NewPhonenumber().GetPhoneNumber(&phonenumber.GetPhoneNumberRequest{
		Code: code,
	})
	if err != nil {
		return err
	}
	if phoneResult.ErrCode != 0 {
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
		rest, err := wxapp.NewClient(app.Weapp).GetUserRiskRank(userRiskRankParam)
		if err != nil {
			app.Logger.Info("BindPhoneByCode 风险等级查询查询出错", err.Error())
		}
		userInfo.Risk = rest.RiskRank
	}

	//获取用户地址  todo 加入队列
	city, err := baidu.IpToCity(cip)
	if err != nil {
		app.Logger.Info("BindPhoneByCode ip地址查询失败", err.Error())
	}
	userInfo.CityCode = city.Content.AddressDetail.Adcode
	userInfo.Ip = cip
	ret := u.r.Save(&userInfo)

	if invitedBy != "" && userInfo.Risk > 2 {
		return errors.New("很遗憾您暂无法参与活动")
	}
	//有邀请，并且没有发放奖励，不是黑产用户，给用户发放奖励
	inviteInfo := u.rInvite.GetInvite(userInfo.OpenId)
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
			app.Logger.Error("发放邀请积分失败", err)
		}
	}
	return ret
}
func (u UserService) BindPhoneByIV(param BindPhoneByIVParam) error {
	userInfo := u.r.GetUserById(param.UserId)
	if userInfo.ID == 0 {
		return errors.New("未查到用户信息")
	}

	session, err := DefaultSessionService.FindSessionByOpenId(userInfo.OpenId)

	if err != nil {
		return err
	}
	if session.ID == 0 {
		return errno.ErrAuth
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
	rest, err := wxapp.NewClient(app.Weapp).GetUserRiskRank(param)
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
	}
	if param.PositionIcon != "" {
		user.PositionIcon = param.PositionIcon
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
	}
	return u.r.Save(&user)
}
