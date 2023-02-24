package weixinmobile

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"time"
)

type RedisAccessTokenCenterConf struct {
	AppId       string
	AppSecret   string
	RedisClient *redis.Client
	Prefix      string
}
type RedisAccessTokenCenter struct {
	RedisAccessTokenCenterConf
	openId string
}

func NewRedisAccessTokenCenterFromOpenId(c RedisAccessTokenCenterConf, openId string) *RedisAccessTokenCenter {
	return &RedisAccessTokenCenter{RedisAccessTokenCenterConf: c, openId: openId}
}

type UserInfo struct {
}

func NewRedisAccessTokenCenterFromCode(c RedisAccessTokenCenterConf, code string) (center *RedisAccessTokenCenter, openId string, err error) {
	resp, err := Code2AccessToken(c.AppId, c.AppSecret, code)
	if err != nil {
		return nil, "", err
	}
	if !resp.IsSuccess() {
		return nil, "", errors.New(fmt.Sprintf("code:%d message:%s", resp.Errcode, resp.Errmsg))
	}
	client := NewRedisAccessTokenCenterFromOpenId(c, resp.Openid)
	client.RedisClient.SetEX(context.Background(), client.refreshAccessTokenKey(), resp.RefreshToken, 30*time.Hour*24)
	client.RedisClient.SetEX(context.Background(), client.accessTokenKey(), resp.AccessToken, time.Duration(resp.ExpiresIn)*time.Second)
	return client, resp.Openid, nil
}

var (
	RefreshTokenExpire = errors.New("refreshtoken expire")
)

func (r RedisAccessTokenCenter) formatKey(key string) string {
	return r.Prefix + key
}
func (r RedisAccessTokenCenter) accessTokenKey() string {
	return r.formatKey("accesstoken:" + r.openId)
}
func (r RedisAccessTokenCenter) refreshAccessTokenKey() string {
	return r.formatKey("refreshtoken:" + r.openId)
}
func (r RedisAccessTokenCenter) AccessToken() (token string, expireIn time.Time, err error) {
	token, err = r.RedisClient.Get(context.Background(), r.accessTokenKey()).Result()
	if err != nil && err != redis.Nil {
		return "", time.Time{}, err
	}

	if token == "" {
		return r.refreshToken("")
	}

	t, err := r.RedisClient.PTTL(context.Background(), r.accessTokenKey()).Result()
	if err != nil || err != redis.Nil {
		return "", time.Time{}, err
	}
	return token, time.Now().Add(t), nil
}

func (r RedisAccessTokenCenter) refreshToken(oldToken string) (token string, expireIn time.Time, err error) {
	refreshToken, err := r.RedisClient.Get(context.Background(), r.refreshAccessTokenKey()).Result()
	if err != nil && err != redis.Nil {
		return "", time.Time{}, err
	}
	if refreshToken == "" {
		return "", time.Time{}, RefreshTokenExpire
	}
	refreshResp, err := RefreshAccessToken(r.AppId, refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	if refreshResp.IsSuccess() {
		r.RedisClient.SetEX(context.Background(), r.refreshAccessTokenKey(), refreshResp.RefreshToken, 30*time.Hour*24)
		r.RedisClient.SetEX(context.Background(), r.accessTokenKey(), refreshResp.AccessToken, time.Duration(refreshResp.ExpiresIn)*time.Second)
		return refreshResp.AccessToken, time.Now().Add(time.Duration(refreshResp.ExpiresIn) * time.Second), nil
	}
	if r.isRefreshTokenExpireCode(refreshResp.Errcode) {
		return "", time.Time{}, RefreshTokenExpire
	}
	return "", time.Time{}, errors.New(fmt.Sprintf("code:%d,message:%s", refreshResp.Errcode, refreshResp.Errmsg))
}
func (r RedisAccessTokenCenter) isRefreshTokenExpireCode(errCode int) bool {
	switch errCode {
	case 40030, 41003, 42002, 42007, 61023:
		return true
	}
	return false
}
func (r RedisAccessTokenCenter) RefreshToken(oldToken string) (token string, expireIn time.Time, err error) {
	return r.refreshToken(oldToken)
}

func (r RedisAccessTokenCenter) IsExpired(accessToken string) (bool, error) {
	resp, err := CheckAccessToken(accessToken, r.openId)
	if err != nil {
		return false, err
	}
	return !resp.IsSuccess(), nil
}

const (
	code2AccessTokenUrl   = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	checkAccessTokenUrl   = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

type Code2AccessTokenResp struct {
	errorResp
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

func (e errorResp) IsSuccess() bool {
	return e.Errcode == 0
}

// Code2AccessToken code换取accesstoken
func Code2AccessToken(appId, secret, code string) (*Code2AccessTokenResp, error) {
	u := fmt.Sprintf(code2AccessTokenUrl, appId, secret, code)
	body, err := httptool.Get(u)
	if err != nil {
		return nil, err
	}
	resp := Code2AccessTokenResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RefreshTokenResp struct {
	errorResp
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

// RefreshAccessToken 刷新accesstoken
func RefreshAccessToken(appId, refreshToken string) (*RefreshTokenResp, error) {
	u := fmt.Sprintf(refreshAccessTokenUrl, appId, refreshToken)
	body, err := httptool.Get(u)
	if err != nil {
		return nil, err
	}
	resp := RefreshTokenResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CheckAccessTokenResp struct {
	errorResp
}

//CheckAccessToken 检查access_token是否有效
func CheckAccessToken(accessToken, openId string) (*CheckAccessTokenResp, error) {
	u := fmt.Sprintf(checkAccessTokenUrl, accessToken, openId)
	body, err := httptool.Get(u)
	if err != nil {
		return nil, err
	}
	resp := CheckAccessTokenResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
