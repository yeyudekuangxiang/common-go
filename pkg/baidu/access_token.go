package baidu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/httputil"
	"time"
)

const TokenUrl = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"

type AccessTokenConfig struct {
	RedisClient *redis.Client
	Prefix      string
	AppKey      string
	AppSecret   string
}

func NewAccessToken(config AccessTokenConfig) *AccessToken {
	lock := util.RedisDistributedLock{Redis: config.RedisClient, Prefix: config.Prefix + "lock:"}
	return &AccessToken{
		cache:     NewRedisCache(config.RedisClient, config.Prefix+"cache:"),
		AppKey:    config.AppKey,
		AppSecret: config.AppSecret,
		lock:      &lock,
	}
}

type AccessToken struct {
	cache     ICache
	AppKey    string
	AppSecret string
	lock      util.DistributedLock
}

type AccessTokenResponse struct {
	ErrorResponse
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	SessionSecret string `json:"session_secret"`
}

func (t *AccessToken) GetToken() (string, error) {
	t.lock.LockWait(fmt.Sprintf("BaiduGetToken%s", t.AppKey), time.Second*50)
	defer t.lock.UnLock(fmt.Sprintf("BaiduGetToken%s", t.AppKey))

	token, err := t.cache.GetValue("accessToken:" + t.AppKey)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}
	token, expireIn, err := t.getToken()
	if err != nil {
		return "", err
	}
	err = t.cache.SetValue("accessToken:"+t.AppKey, token, expireIn-time.Second*300)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (t *AccessToken) getToken() (string, time.Duration, error) {
	u := fmt.Sprintf(TokenUrl, t.AppKey, t.AppSecret)
	body, err := httputil.PostMapFrom(u, map[string]string{})
	if err != nil {
		return "", 0, err
	}
	resp := AccessTokenResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", 0, err
	}
	if !resp.IsSuccess() {
		return "", 0, errors.New(resp.ErrorDescription)
	}
	return resp.AccessToken, time.Second * time.Duration(resp.ExpiresIn), nil
}
func (t *AccessToken) RefreshToken() (string, error) {
	t.lock.LockWait(fmt.Sprintf("BaiduGetToken%s", t.AppKey), time.Second*5)
	defer t.lock.UnLock(fmt.Sprintf("BaiduGetToken%s", t.AppKey))

	token, expireIn, err := t.getToken()
	if err != nil {
		return "", err
	}
	err = t.cache.SetValue("accessToken:"+t.AppKey, token, expireIn-time.Second*300)
	if err != nil {
		return "", err
	}
	return token, nil
}
