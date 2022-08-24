package oola

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/duiba/util"
	"net/url"
)

type Oola struct {
	ctx        *context.MioContext
	appId      string
	clientId   string
	oolaUserId string
	sign       string
	domain     string
	redis      *redis.Client
}

type oolaResponse struct {
	Status string       `json:"status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Msg    string       `json:"msg,omitempty"`
	Info   responseInfo `json:"info,omitempty"`
}

type responseInfo struct {
	OolaUserId   string `json:"oolaUserId,omitempty"`
	AutologinKey string `json:"autologinKey,omitempty"`
}

func NewOola(context *context.MioContext, appId, clientId, sign, domain string, client *redis.Client) *Oola {
	return &Oola{
		ctx:      context,
		appId:    appId,
		clientId: clientId,
		sign:     sign,
		domain:   domain,
		redis:    client,
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

func (o Oola) GetToken() (string, error) {
	autoLoginKey, err := o.redis.Get(o.ctx, "oola_login_key:"+o.clientId).Result()
	if err != nil {
		if err == redis.Nil {
			//重新获取token
			_, oolaUserLoginKey, err := o.register()
			if err != nil {
				return "", err
			}
			return oolaUserLoginKey, nil
		}
		app.Logger.Error(err)
		return "", err
	}
	if autoLoginKey == "" {
		return "", errors.New("数据异常")
	}
	return autoLoginKey, nil
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
	//记录redis todo
	return res.Info.OolaUserId, res.Info.AutologinKey, nil
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
	//记录redis todo
	return res.Info.OolaUserId, res.Info.AutologinKey, nil
}
