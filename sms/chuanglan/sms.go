package chuanglan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SmsClient struct {
	Account  string
	Password string
	BaseUrl  string
}

func NewSmsClient(account string, password string, BaseUrl string) *SmsClient {
	return &SmsClient{
		Account:  account,
		Password: password,
		BaseUrl:  BaseUrl,
	}
}

type SmsReturn struct {
	Code       string `json:"code"`
	FailNum    string `json:"failNum"`
	SuccessNum string `json:"successNum"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
	ErrorMsg   string `json:"errorMsg"`
}

func (c *SmsClient) Send(mobile string, content string) (*SmsReturn, error) {
	url := c.BaseUrl
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + c.Account + `",
    "password": "` + c.Password + `", //需要加入K8S
    "msg": "` + content + `",
    "phone": "` + mobile + `",
    "sendtime": "201704101400",
    "report": "true",
    "extend": "555",
    "uid": "321abc"
}`)
	fmt.Println(payload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
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
