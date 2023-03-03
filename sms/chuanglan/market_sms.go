package chuanglan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MarketSmsClient struct {
	Account  string
	Password string
	BaseUrl  string
}

func NewMarketSmsClient(account string, password string, baseUrl string) *MarketSmsClient {
	return &MarketSmsClient{
		Account:  account,
		Password: password,
		BaseUrl:  baseUrl,
	}
}

type MarketSmsReturn struct {
	Code       string `json:"code"`
	FailNum    string `json:"failNum"`
	SuccessNum string `json:"successNum"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
	ErrorMsg   string `json:"errorMsg"`
}

func (c *MarketSmsClient) Send(mobile string, content string) (*SmsReturn, error) {
	params := make(map[string]string, 0)
	//请登录zz.253.com获取API账号、密码以及短信发送的URL
	params["account"] = c.Account   //创蓝API账号
	params["password"] = c.Password //创蓝API密码
	params["phone"] = mobile
	params["msg"] = "【绿喵mio】" + content
	params["report"] = "true"
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", c.BaseUrl, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := &SmsReturn{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	if ret.Code != "0" {
		return ret, fmt.Errorf("send sms failed, response code: %s, msgId: %s", ret.Code, ret.MsgId)
	}
	return ret, nil
}
