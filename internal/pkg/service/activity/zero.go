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
	"time"
)

var (
	ZeroActivityStartTime, _ = time.Parse("2006-01-02 15:04:05", "2022-04-11 14:08:00")
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
	path := "https://88543.activity-12.m.duiba.com.cn/aaw/haggle/index?opId=194935804526281&dbnewopen&newChannelType=3"
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
		DCustom: fmt.Sprintf("avatar=%s&nickname=%snewUser=%d", userInfo.AvatarUrl, userInfo.Nickname, isNewUser),
	})
}
func (srv ZeroService) StoreUrl(url string) (string, error) {
	key := util.Md5(url)
	redisKey := fmt.Sprintf(config.RedisKey.DuiBaShortUrl, key)
	err := app.Redis.Set(context.Background(), redisKey, url, time.Hour*7).Err()
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
