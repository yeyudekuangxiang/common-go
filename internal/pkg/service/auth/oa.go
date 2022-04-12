package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chanxuehong/wechat/mp/jssdk"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	wxoa2 "mio/pkg/wxoa"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	LoginUrl = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
)

type OaService struct {
	Platform entity.UserSource
}

func (srv OaService) LoginByCode(code string) (string, error) {
	setting := config.FindOaSetting(srv.Platform)

	oauth2Client := oauth2.Client{
		Endpoint: mpoauth2.NewEndpoint(setting.AppId, setting.Secret),
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		return "", nil
	}

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		return "", nil
	}
	sexStr := ""
	if userinfo.Sex == 1 {
		sexStr = " MALE"
	} else if userinfo.Sex == 2 {
		sexStr = "FEMALE"
	}
	user, err := service.DefaultUserService.CreateUser(service.CreateUserParam{
		OpenId:    userinfo.OpenId,
		AvatarUrl: userinfo.HeadImageURL,
		Gender:    sexStr,
		Nickname:  userinfo.Nickname,
		Source:    srv.Platform,
		UnionId:   userinfo.UnionId,
	})
	if err != nil {
		return "", err
	}
	return service.DefaultUserService.CreateUserToken(user.ID)
}
func (srv OaService) CheckAuthWhiteList(platform entity.UserSource, url string) bool {
	return true
}
func (srv OaService) AutoLoginCallback(code string, state string) (string, error) {

	redisKey := fmt.Sprintf("%s%s", config.RedisKey.OaAuth, state)
	dataStr, err := app.Redis.Get(context.Background(), redisKey).Result()
	if err != nil {
		app.Logger.Error(err)
	}
	if dataStr == "" {
		return "", errors.New("数据异常")
	}
	data := map[string]string{}
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return "", err
	}
	var redirectUri string
	if index := strings.Index(data["RedirectUri"], "?"); index >= 0 {
		prefix := data["RedirectUri"][:index+1]
		last := data["RedirectUri"][index+1:]
		redirectUri = fmt.Sprintf("%scode=%s&state=%s&platform=%s&%s", prefix, code, state, data["App"], last)
	} else if index := strings.Index(data["RedirectUri"], "#"); index >= 0 {
		prefix := data["RedirectUri"][:index]
		last := data["RedirectUri"][index:]
		redirectUri = fmt.Sprintf("%s?code=%s&state=%s&platform=%s%s", prefix, code, state, data["App"], last)
	} else {
		redirectUri = fmt.Sprintf("%s?code=%s&state=%s&platform=%s", data["RedirectUri"], code, state, data["App"])
	}

	app.Redis.Del(context.Background(), redisKey)

	app.Logger.Info("授权回调:", redirectUri)
	//如果使用302重定向 浏览器点击返回上一页时 url会变成此方法url且会重新进入此方法报数据异常错误
	//使用301重定向时 浏览器点击返回上一页时 url会变成此方法url但不会进入此方法 会已相同的code进入回调页
	//c.Redirect(http.StatusMovedPermanently, redirectUri)
	return redirectUri, nil
}
func (srv OaService) AutoLogin(redirectUri string, state string) (string, error) {
	if !srv.CheckAuthWhiteList(srv.Platform, redirectUri) {
		return "", errors.New("跳转地址未在白名单内")
	}

	setting := config.FindOaSetting(srv.Platform)

	data := map[string]string{
		"RedirectUri": redirectUri,
		"State":       state,
		"App":         string(srv.Platform),
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	t := time.Now().UnixNano()
	key := util.Md5(fmt.Sprintf("%s%d", setting.AppId, t))
	redisKey := fmt.Sprintf("%s%s", config.RedisKey.OaAuth, key)

	err = app.Redis.Set(context.Background(), redisKey, string(dataBytes), 30*time.Second).Err()
	if err != nil {
		return "", err
	}
	jumpUrl := fmt.Sprintf("%s%s", config.Config.App.Domain, "/api/mp2c/oa/auth/callback")
	escapeUrl := url.QueryEscape(jumpUrl)
	loginUrl := fmt.Sprintf(LoginUrl, setting.AppId, escapeUrl, key)
	app.Logger.Info("跳转登陆url", loginUrl)
	//使用301重定项时 后续的授权将会已相同的state直接进入callback函数 而不会经过本函数
	//使用302重定向时 后续的授权将会同样先进入本函数 保存数据后进行授权然后进入callback函数
	//c.Redirect(http.StatusFound, loginUrl)
	return loginUrl, nil
}
func (srv OaService) Sign(url string) (*OaSignResult, error) {
	setting := config.FindOaSetting(srv.Platform)

	tickerServer := wxoa2.TicketTokenServer{
		TokenServer: &wxoa2.AccessTokenServer{
			AppId:  setting.AppId,
			Secret: setting.Secret,
			Redis:  app.Redis,
		},
		AppId: setting.AppId,
		Redis: app.Redis,
	}

	ticker, err := tickerServer.Ticket()
	if err != nil {
		app.Logger.Error(err)
		return nil, errno.InternalServerError
	}

	nonceStr := util.Md5(time.Now().String())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := jssdk.WXConfigSign(ticker, nonceStr, timestamp, url)
	return &OaSignResult{
		AppId:     setting.AppId,
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Signature: sign,
	}, nil
}
