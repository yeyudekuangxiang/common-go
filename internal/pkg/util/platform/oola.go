package platform

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util"
	"net/url"
	"sort"
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

func (o *Oola) WithPhone(phone string) {
	if phone != "" {
		o.phone = phone
	}
}

func (o *Oola) WithHeadImgUrl(headImgUrl string) {
	if headImgUrl != "" {
		o.headImgUrl = headImgUrl
	}
}

func (o *Oola) WithUserName(userName string) {
	if userName != "" {
		o.userName = userName
	}
}

func (o *Oola) getSign(key string, params url.Values) (sign string, err error) {
	var signStr string
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	for _, v := range slice {
		signStr += v + "=" + util.InterfaceToString(params[v][0]) + ";"
	}
	return encrypttool.Md5(key + signStr), nil
}

func (o *Oola) GetToken(key string) (string, string, error) {
	autoLoginKey, _ := o.redis.GetDel(o.ctx, "oola_login_key:"+o.clientId).Result()
	channelCode, _ := o.redis.GetDel(o.ctx, "oola_channel_code:"+o.clientId).Result()
	if autoLoginKey == "" || channelCode == "" {
		return o.register(key)
	}
	return channelCode, autoLoginKey, nil
}

func (o *Oola) register(key string) (string, string, error) {
	params := make(url.Values)
	params.Set("appId", o.appId)
	params.Set("clientId", o.clientId)
	if o.userName != "" {
		params.Set("userName", o.userName)
	}
	if o.phone != "" {
		params.Set("phone", o.phone)
	}
	if o.headImgUrl != "" {
		params.Set("headImgUrl", o.headImgUrl)
	}

	sign, err := o.getSign(key, params)
	if err != nil {
		return "", "", err
	}
	params.Set("sign", sign)
	u := o.domain + "/api/user/register"
	body, err := httptool.PostFrom(u, params)
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
			return o.getUserAutoLoginKey(key)
		}
		return "", "", errors.New(res.Msg)
	}
	//记录redis
	o.redis.SetNX(o.ctx, "oola_login_key:"+o.clientId, res.Info.AutologinKey, 9*time.Minute)
	o.redis.SetNX(o.ctx, "oola_channel_code:"+o.clientId, res.Info.ChannelCode, 9*time.Minute)
	return res.Info.ChannelCode, res.Info.AutologinKey, nil
}

func (o *Oola) getUserAutoLoginKey(key string) (string, string, error) {
	params := make(url.Values)
	params.Set("appId", o.appId)
	params.Set("clientId", o.clientId)
	sign, err := o.getSign(key, params)
	if err != nil {
		return "", "", err
	}
	params.Set("sign", sign)
	u := o.domain + "/api/user/getUserAutoLoginKey"
	body, err := httptool.PostFrom(u, params)
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
	o.redis.SetNX(o.ctx, "oola_login_key:"+o.clientId, res.Info.AutologinKey, 9*time.Minute)
	o.redis.SetNX(o.ctx, "oola_channel_code:"+o.clientId, res.Info.ChannelCode, 9*time.Minute)
	return res.Info.ChannelCode, res.Info.AutologinKey, nil
}
