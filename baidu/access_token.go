package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/yeyudekuangxiang/common-go/tool/httptool"
	"time"
)

const TokenUrl = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"

type AccessTokenConfig struct {
	Cache     ICache
	AppKey    string
	AppSecret string
}

func NewAccessToken(config AccessTokenConfig) *AccessToken {
	return &AccessToken{
		cache:     config.Cache,
		AppKey:    config.AppKey,
		AppSecret: config.AppSecret,
	}
}

type RedisAccessTokenConfig struct {
	RedisClient *redis.Client
	Prefix      string
	AppKey      string
	AppSecret   string
}

func NewRedisAccessToken(config RedisAccessTokenConfig) *AccessToken {
	return &AccessToken{
		cache:     NewRedisCache(config.RedisClient, config.Prefix+config.AppKey+":"),
		AppKey:    config.AppKey,
		AppSecret: config.AppSecret,
	}
}

type MemoryAccessTokenConfig struct {
	Prefix    string
	AppKey    string
	AppSecret string
}

func NewMemoryAccessToken(config MemoryAccessTokenConfig) *AccessToken {
	return &AccessToken{
		cache:     NewMemoryCache(config.Prefix + config.AppKey + ":"),
		AppKey:    config.AppKey,
		AppSecret: config.AppSecret,
	}
}

type IAccessToken interface {
	GetToken() (string, error)
	RefreshToken() (string, error)
}
type AccessToken struct {
	cache     ICache
	AppKey    string
	AppSecret string
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
	body, err := httptool.PostMapFrom(u, map[string]string{})
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
