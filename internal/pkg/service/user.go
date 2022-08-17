package service

import (
	"context"
	"fmt"
	"github.com/medivhzhan/weapp/v3/phonenumber"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	util2 "mio/internal/pkg/util"
	"mio/internal/pkg/util/message"
	"mio/pkg/errno"
	"mio/pkg/wxapp"
	"strconv"
	"time"
)

var DefaultUserService = NewUserService(repository2.DefaultUserRepository)

func NewUserService(r repository2.IUserRepository) UserService {
	return UserService{
		r: r,
	}
}

type UserService struct {
	r repository2.IUserRepository
}

func (u UserService) GetUserById(id int64) (*entity.User, error) {
	if id == 0 {
		return &entity.User{}, nil
	}
	user := u.r.GetUserById(id)
	return &user, nil
}
func (u UserService) GetUserByOpenId(openId string) (*entity.User, error) {
	if openId == "" {
		return &entity.User{}, nil
	}
	user := u.r.GetUserBy(repository2.GetUserBy{
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

func (u UserService) CreateUser(param CreateUserParam) (*entity.User, error) {
	user := u.r.GetUserBy(repository2.GetUserBy{
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
	return &user, repository2.DefaultUserRepository.Save(&user)
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
func (u UserService) GetUserBy(by repository2.GetUserBy) (*entity.User, error) {
	user := repository2.DefaultUserRepository.GetUserBy(by)
	return &user, nil
}
func (u UserService) FindOrCreateByMobile(mobile string, cid int64) (*entity.User, error) {
	if mobile == "" {
		return nil, errors.New("手机号不能为空")
	}
	user := repository2.DefaultUserRepository.GetUserBy(repository2.GetUserBy{
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

// FindUserBySource 根据用户id 获取指定平台的用户
func (u UserService) FindUserBySource(source entity.UserSource, userId int64) (*entity.User, error) {
	if userId == 0 {
		return &entity.User{}, nil
	}

	user := repository2.DefaultUserRepository.GetUserById(userId)

	if user.ID == 0 || user.PhoneNumber == "" {
		return &entity.User{}, nil
	}

	sourceUer := repository2.DefaultUserRepository.GetUserBy(repository2.GetUserBy{
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
func (u UserService) BindPhoneByCode(userId int64, code string) error {
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
	return u.r.Save(&userInfo)
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
	list := pointTranService.GetListBy(repository2.GetPointTransactionListBy{
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
func (u UserService) GetUserListBy(by repository2.GetUserListBy) ([]entity.User, error) {
	return u.r.GetUserListBy(by), nil
}
func (u UserService) UpdateUserInfo(param UpdateUserInfoParam) error {
	user := u.r.GetUserById(param.UserId)
	if user.ID == 0 {
		return errno.ErrUserNotFound
	}
	user.AvatarUrl = param.Avatar
	user.Nickname = param.Nickname
	user.Gender = param.Gender
	return u.r.Save(&user)
}

func (u UserService) GetUserPageListBy(by repository2.GetUserPageListBy) ([]entity.User, int64) {
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
	certCount, err := DefaultBadgeService.GetUserCertCountById(userId)
	if err != nil {
		return nil, err
	}
	return &UserAccountInfo{
		Balance:     point.Balance,
		CertNum:     certCount,
		CarbonToday: "100g",
		CarbonAll:   "101g",
	}, nil
}

/**更新carbon*/

func (u UserService) UpdateUserCarbon(uid int64, value float64) {
	valDec := decimal.NewFromFloat(value)
	if valDec.IsZero() || uid == 0 {
		return
	}
	user := u.r.GetUserById(uid)
	if user.ID == 0 {
		return
	}
	carbonDec := decimal.NewFromFloat(user.Carbon)
	addVal, addErr := carbonDec.Add(valDec).Float64()
	if !addErr {
		return
	}
	user.Carbon = addVal
	err := u.r.Save(&user)
	if err != nil {
		app.Logger.Error("user表更新carbon失败", uid, value, err)
	}
}
