package chuanglan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type MarketTemplateSmsClient struct {
	Account  string
	Password string
	BaseUrl  string
}

func NewMarketTemplateSmsClient(account string, password string, baseUrl string) *MarketTemplateSmsClient {
	return &MarketTemplateSmsClient{
		Account:  account,
		Password: password,
		BaseUrl:  baseUrl,
	}
}

type MarketTemplateSmsReturn struct {
	Code       string `json:"code"`
	FailNum    string `json:"failNum"`
	SuccessNum string `json:"successNum"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
	ErrorMsg   string `json:"errorMsg"`
}

func (c *MarketTemplateSmsClient) Send(phone string, templateContent string, msg string) (smsReturn *SmsReturn, err error) {
	params := ""
	if msg == "" {
		params = phone //组装 18840853003,小李,1;
	} else {
		params = phone + "," + msg + ";" //组装 18840853003,小李,1;
	}

	url := c.BaseUrl
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + c.Account + `",
    "password": "` + c.Password + `", //需要加入K8S
    "msg": "` + templateContent + `",
	"params":"` + params + `",
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
