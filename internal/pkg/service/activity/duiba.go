package activity

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"mio/pkg/wxapp"
	"time"
)

var (
	ZeroActivityStartTime, _ = time.Parse("2006-01-02 15:04:05", "2022-04-13 00:00:00")
)
var DefaultZeroService = ZeroService{}

type ZeroService struct {
}

func (srv ZeroService) AutoLogin(userId int64, short string) (string, error) {
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return "", err
	}
	if userInfo.ID == 0 {
		return "", errors.New("未查询到用户信息")
	}

	//此方法只能使用一次
	isNewUser, err := srv.IsNewUser(userId, userInfo.Time.Time)
	if err != nil {
		return "", err
	}
	path := "https://88543.activity-12.m.duiba.com.cn/aaw/haggle/index?opId=195380425492081&dbnewopen&newChannelType=3"
	if short != "" {
		p, err := srv.GetUrlByShort(short)
		if err != nil {
			app.Logger.Error(userId, short, err)
		}
		if p != "" {
			path = p
		}
	}

	return service.DefaultDuiBaService.AutoLoginOpenId(service.AutoLoginOpenIdParam{
		UserId:  userId,
		OpenId:  userInfo.OpenId,
		Path:    path,
		DCustom: fmt.Sprintf("avatar=%s&nickname=%s&newUser=%d", userInfo.AvatarUrl, userInfo.Nickname, isNewUser),
	})
}
func (srv ZeroService) StoreUrl(url string) (string, error) {
	key := util.Md5(url)
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, key)
	err := app.Redis.Set(context.Background(), redisKey, url, time.Hour*10*24).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}
func (srv ZeroService) GetUrlByShort(short string) (string, error) {
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, short)
	u, err := app.Redis.Get(context.Background(), redisKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return u, nil
}
func (srv ZeroService) IsNewUser(userId int64, createTime time.Time) (int, error) {
	//用户创建时间在活动开始时间之前
	if createTime.Before(ZeroActivityStartTime) {
		return 0, nil
	}
	return 1, nil

	//兑吧自己有判断 不用我们记录了
	//判断是否已经助力过
	redisKey := fmt.Sprintf(config.RedisKey.ActivityZeroIsNewUser, userId)
	result, err := app.Redis.SetNX(context.Background(), redisKey, time.Now().Unix(), time.Hour*24*30).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return util.Ternary(result, 1, 0).Int(), nil
}

type DuiBaActivity struct {
	ActivityId  string
	ActivityUrl string
	StartTime   time.Time
	EndTime     time.Time
}

const DUIBAIndex = "https://88543.activity-12.m.duiba.com.cn/chw/visual-editor/skins?id=239935"

func (srv ZeroService) DuiBaStoreUrl(activityId string, url string) (string, error) {
	key := activityId + "_" + util.Md5(url)
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, key)
	err := app.Redis.Set(context.Background(), redisKey, url, time.Hour*10*24).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}
func (srv ZeroService) GetDuiBaUrlByShort(short string) (string, error) {
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, short)
	u, err := app.Redis.Get(context.Background(), redisKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return u, nil
}
func (srv ZeroService) DuiBaAutoLogin(userId int64, activityId, short, thirdParty string, cip string) (string, error) {
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return "", err
	}
	if userInfo.ID == 0 {
		return "", errno.ErrUserNotFound
	}

	path := DUIBAIndex
	isNewUser := false

	activity, err := service.DefaultDuiBaActivityService.FindActivity(activityId)
	if err != nil {
		return "", err
	}

	if activity.ID != 0 {
		isNewUser, err = srv.IsDuiBaActivityNewUser(activityId, userId)
		if err != nil {
			return "", err
		}

		path = activity.ActivityUrl

		if short != "" {
			p, err := srv.GetDuiBaUrlByShort(short)
			if err != nil {
				app.Logger.Error(userId, short, err)
			}
			if p != "" {
				path = p
			}
		}
	}

	//检测用户风险等级
	userRiskRankParam := wxapp.UserRiskRankParam{
		AppId:    config.Config.Weapp.AppId,
		OpenId:   userInfo.OpenId,
		Scene:    1,
		ClientIp: cip,
	}
	//风险等级4,代表不需要验证.0-3代表需要做出限制
	if activity.RiskLimit < 4 {
		//判断用户手机号,警告:必须是非首页
		if userInfo.PhoneNumber == "" {
			return "", errno.ErrNotBindMobile
		}
		userRiskRankParam.MobileNo = userInfo.PhoneNumber
	}
	resp, err := service.DefaultUserService.CheckUserRisk(userRiskRankParam)
	if err != nil && activityId != "index" {
		app.Logger.Info("DuiBaAutoLogin 风险等级查询查询出错", err.Error())
	} else {
		if resp.RiskRank > activity.RiskLimit {
			return "", nil
			//return "", errors.New("该活动仅限部分地区用户参与,请看看其他活动")
		}
	}

	vip := 0
	if thirdParty == "thirdParty" && isNewUser {
		vip = 2
	}

	isNewUserInt := util.Ternary(isNewUser, 1, 0).Int()
	return service.DefaultDuiBaService.AutoLoginOpenId(service.AutoLoginOpenIdParam{
		UserId:  userId,
		OpenId:  userInfo.OpenId,
		Path:    path,
		Vip:     vip,
		DCustom: fmt.Sprintf("avatar=%s&nickname=%s&newUser=%d", userInfo.AvatarUrl, userInfo.Nickname, isNewUserInt),
	})
}
func (srv ZeroService) IsDuiBaActivityNewUser(activityId string, userId int64) (bool, error) {
	userInfo, err := service.DefaultUserService.GetUserById(userId)
	if err != nil {
		return false, err
	}
	if userInfo.ID == 0 {
		return false, errno.ErrUserNotFound
	}
	return !userInfo.Time.IsZero() && userInfo.Time.Add(time.Hour*24).After(time.Now()), nil
}
