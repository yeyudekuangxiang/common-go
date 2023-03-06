package chuanglan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	queryMarket         = "/msg/v1/send/json"
	queryMarketTemplate = "/msg/variable/json"
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
	errorResp
	Code       string `json:"code"`
	FailNum    string `json:"failNum"`
	SuccessNum string `json:"successNum"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
	ErrorMsg   string `json:"errorMsg"`
}

func (c *MarketSmsClient) Send(mobile string, content string, sign string) (*SmsReturn, error) {
	url := domain + queryMarket
	params := make(map[string]string, 0)
	//请登录zz.253.com获取API账号、密码以及短信发送的URL
	params["account"] = c.Account   //创蓝API账号
	params["password"] = c.Password //创蓝API密码
	params["phone"] = mobile
	params["msg"] = sign + content
	params["report"] = "true"
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	res, err := client.Do(request)
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
	return ret, err
}

func (c *MarketSmsClient) SendTemplateSms(phone string, templateContent string, sign string, paramsSlice []string) (smsReturn *SmsReturn, err error) {
	url := domain + queryMarketTemplate
	params := ""
	if len(paramsSlice) == 0 {
		params = phone //组装 18840853003,小李,1;
	} else {
		params = phone + "," + strings.Join(paramsSlice, ",") + ";" //组装 18840853003,小李,1;
	}
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + c.Account + `",
    "password": "` + c.Password + `", //需要加入K8S
    "msg": "` + sign + templateContent + `",
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
	return ret, err
}
