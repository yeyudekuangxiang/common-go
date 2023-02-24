package weixinmobile

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	userInfoUrl = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
)

type Lang string

const (
	ZhCN Lang = "zh_CN"
	ZhTw Lang = "zh_TW"
	EN   Lang = "en"
)

type UserInfoResp struct {
	errorResp
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

//GetUserInfo 获取用户信息
func GetUserInfo(accessToken, openId string, lang Lang) (*UserInfoResp, error) {
	u := fmt.Sprintf(userInfoUrl, accessToken, openId)
	body, err := httptool.Get(u)
	if err != nil {
		return nil, err
	}
	resp := UserInfoResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
