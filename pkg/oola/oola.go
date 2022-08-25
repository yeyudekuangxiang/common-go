package oola

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/duiba/util"
	"net/url"
	"time"
)

type Oola struct {
	ctx        *context.MioContext
	appId      string
	clientId   string
	oolaUserId string
	domain     string
	redis      *redis.Client
	headImgUrl string
	userName   string
	phone      string
}

type oolaResponse struct {
	Status string       `json:"status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Msg    string       `json:"msg,omitempty"`
	Info   responseInfo `json:"info,omitempty"`
}

type responseInfo struct {
	OolaUserId   string `json:"oolaUserId"`
	AutologinKey string `json:"autologinKey"`
	ChannelCode  string `json:"channelCode"`
}

func NewOola(context *context.MioContext, appId, clientId, domain string, client *redis.Client) *Oola {
	return &Oola{
		ctx:      context,
		appId:    appId,
		clientId: clientId,
		domain:   domain,
		redis:    client,
	}
}

func (o Oola) WithPhone(phone string) {
	if phone != "" {
		o.phone = phone
	}
}

func (o Oola) WithHeadImgUrl(headImgUrl string) {
	if headImgUrl != "" {
		o.headImgUrl = headImgUrl
	}
}

func (o Oola) WithUserName(userName string) {
	if userName != "" {
		o.userName = userName
	}
}

func (o Oola) getSign(ch string) (sign string, err error) {
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(ch)
	if scene.Key == "" || scene.Key == "e" {
		return "", errors.New("渠道查询失败")
	}
	return util.Md5(scene.Key + "appId=" + scene.AppId + ";clientId=" + o.clientId + ";"), nil
}

func (o Oola) GetToken() (string, string, error) {
	autoLoginKey, _ := o.redis.Get(o.ctx, "oola_login_key:"+o.clientId).Result()
	channelCode, _ := o.redis.Get(o.ctx, "oola_channel_code:"+o.clientId).Result()
	if autoLoginKey == "" || channelCode == "" {
		return o.register()
	}
	return channelCode, autoLoginKey, nil
}

func (o Oola) register() (string, string, error) {
	sign, err := o.getSign("oola")
	if err != nil {
		return "", "", err
	}
	params := make(url.Values)
	params.Set("appId", o.appId)
	params.Set("clientId", o.clientId)
	params.Set("sign", sign)

	if o.userName != "" {
		params.Set("userName", o.userName)
	}
	if o.phone != "" {
		params.Set("phone", o.phone)
	}
	if o.headImgUrl != "" {
		params.Set("headImgUrl", o.headImgUrl)
	}

	u := o.domain + "/api/user/register"
	body, err := httputil.PostFrom(u, params)
	if err != nil {
		return "", "", err
	}
	//response
	res := oolaResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", "", err
	}
	if res.Code != "200" {
		if res.Code == "1006" {
			return o.getUserAutoLoginKey()
		}
		return "", "", errors.New(res.Msg)
	}
	//记录redis
	o.redis.Set(o.ctx, "oola_login_key:"+o.clientId, res.Info.AutologinKey, 9*time.Minute)
	o.redis.Set(o.ctx, "oola_channel_code:"+o.clientId, res.Info.ChannelCode, 9*time.Minute)
	return res.Info.ChannelCode, res.Info.AutologinKey, nil
}

func (o Oola) getUserAutoLoginKey() (string, string, error) {
	sign, err := o.getSign("oola")
	if err != nil {
		return "", "", err
	}
	params := make(url.Values)
	params.Set("appId", o.appId)
	params.Set("clientId", o.clientId)
	params.Set("sign", sign)
	u := o.domain + "/api/user/getUserAutoLoginKey"
	body, err := httputil.PostFrom(u, params)
	if err != nil {
		return "", "", err
	}
	//response
	res := oolaResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", "", err
	}
	if res.Code != "200" {
		return "", "", errors.New(res.Msg)
	}
	//记录redis
	o.redis.Set(o.ctx, "oola_login_key:"+o.clientId, res.Info.AutologinKey, 9*time.Minute)
	o.redis.Set(o.ctx, "oola_channel_code:"+o.clientId, res.Info.ChannelCode, 9*time.Minute)
	return res.Info.ChannelCode, res.Info.AutologinKey, nil
}
