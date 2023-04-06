package chuanglan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	domain = "https://smssh1.253.com"
	query  = "/msg/v1/send/json"
)

type SmsClient struct {
	Account  string
	Password string
}

func NewSmsClient(account string, password string) *SmsClient {
	return &SmsClient{
		Account:  account,
		Password: password,
	}
}

type SmsReturn struct {
	ErrorResp
	FailNum    string `json:"failNum"`
	SuccessNum string `json:"successNum"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
}

//发验证码短信 mobile 手机号 content 短信内容  sign 签名，不填默认是【绿喵mio】企业是【企业员工碳减排管理平台】

func (c *SmsClient) Send(mobile string, content string, sign string) (*SmsReturn, error) {
	url := domain + query
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + c.Account + `",
    "password": "` + c.Password + `", //需要加入K8S
    "msg": "` + sign + content + `",
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
	return ret, nil
}