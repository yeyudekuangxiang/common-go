package duiba

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/duiba/util"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"net/url"
	"strconv"
	"time"
)

const (
	baseUrl                 = "https://88543.activity-12.m.duiba.com.cn"
	autoLoginPath           = "/autoLogin/autologin"
	addActivityVistTimesUrl = "/activityVist/addTimes"
)

type Client struct {
	AppKey    string
	AppSecret string
	validator binding.StructValidator
}

func NewClient(appKey, appSecret string) *Client {
	return &Client{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
}
func (client Client) sign(v Param) (map[string]string, error) {
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

	return baseUrl + autoLoginPath + "?" + util.BuildQuery(signParams), nil
}
func (client Client) CheckSign(v Param) error {
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
func (client Client) AddActivityTimes(param AddActivityTimesParam) (*AddActivityTimesParamResp, error) {
	signParams, err := client.sign(param)
	if err != nil {
		return nil, err
	}
	body, err := httptool.Get(baseUrl + addActivityVistTimesUrl + "?" + util.BuildQuery(signParams))
	if err != nil {
		return nil, err
	}
	resp := AddActivityTimesParamResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
