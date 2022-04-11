package duiba

import (
	"github.com/pkg/errors"
	"mio/pkg/validator"
	"net/url"
	"strconv"
	"time"
)

const (
	baseUrl       = "https://activity.m.duiba.com.cn"
	autoLoginPath = "/autoLogin/autologin"
)

type Client struct {
	AppKey    string
	AppSecret string
}

func NewClient(appKey, appSecret string) *Client {
	return &Client{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
}
func (client Client) sign(v Param) (map[string]string, error) {
	err := validator.NewValidator().ValidateStruct(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	params := v.ToMap()
	params["appKey"] = client.AppKey
	params["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	params["appSecret"] = client.AppSecret
	params["sign"] = sign(params)
	delete(params, "appSecret")
	return params, nil
}

// AutoLogin 获取免登陆地址
func (client Client) AutoLogin(param AutoLoginParam) (string, error) {
	signParams, err := client.sign(param)
	if err != nil {
		return "", err
	}
	if dcustom := signParams["dcustom"]; dcustom != "" {
		signParams["dcustom"] = url.QueryEscape(signParams["dcustom"])
	}
	if redirect := signParams["redirect"]; redirect != "" {
		signParams["redirect"] = url.QueryEscape(signParams["redirect"])
	}
	if redirect := signParams["transfer"]; redirect != "" {
		signParams["transfer"] = url.QueryEscape(signParams["transfer"])
	}

	return baseUrl + autoLoginPath + "?" + BuildQuery(signParams), nil
}
func (client Client) CheckSign(v Param) error {
	err := validator.NewValidator().ValidateStruct(v)
	if err != nil {
		return errors.WithStack(err)
	}
	params := v.ToMap()
	params["appKey"] = client.AppKey
	params["appSecret"] = client.AppSecret
	originSign := params["sign"]
	delete(params, "sign")
	if originSign == sign(params) {
		return nil
	}
	return errors.New("签名验证失败")
}
