package message

import (
	"fmt"
	"github.com/square/go-jose/v3/json"
	"io/ioutil"
	"mio/config"
	"mio/pkg/errno"
	"net/http"
	"strings"
)

//{"code":"0","failNum":"0","successNum":"1","msgId":"22110915322300602201000033772693","time":"20221109153223","errorMsg":""}
//{"code":"102","msgId":"","time":"20221109153305","errorMsg":"密码错误"}
//发送短信返回结构

type SmsReturn struct {
	Code       string `json:"code"`
	FailNum    string `json:"failNum"`
	SuccessNum string `json:"successNum"`
	MsgId      string `json:"msgId"`
	Time       string `json:"time"`
	ErrorMsg   string `json:"errorMsg"`
}

//发送验证码

func SendYZMSms(mobile string, code string) (smsReturn *SmsReturn, err error) {
	url := config.Config.Sms.Url
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + config.Config.Sms.Account + `",
    "password": "` + config.Config.Sms.Password + `", //需要加入K8S
    "msg": "验证码` + code + `，30分钟有效。参与低碳任务，体验格调生活。如非本人操作请忽略。",
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
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ret := &SmsReturn{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	if ret.Code != "0" {
		return ret, errno.ErrCommon.WithMessage(ret.ErrorMsg)
	}
	return ret, nil
}

//发送营销短信，也叫模版短信

func SendMarketSms(templateContent string, phone string, msg string) (smsReturn *SmsReturn, err error) {
	msg = phone + "," + msg + ";" //组装 18840853003,小李,1;
	url := config.Config.SmsMarket.Url
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + config.Config.SmsMarket.Account + `",
    "password": "` + config.Config.SmsMarket.Password + `", //需要加入K8S
    "msg": "` + templateContent + `",
	"params":"` + msg + `",
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
	//fmt.Println(string(body))
	ret := &SmsReturn{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	if ret.Code != "0" {
		return ret, errno.ErrCommon.WithMessage(ret.ErrorMsg)
	}
	return ret, nil
}

//验证码：{s}，10分钟有效。如非本人操作，请忽略

func SendYZMSms2B(mobile string, code string) (smsReturn *SmsReturn, err error) {
	url := config.Config.Sms.Url
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + config.Config.Sms.Account + `",
    "password": "` + config.Config.Sms.Password + `", //需要加入K8S
    "msg": "【企业员工碳减排管理平台】验证码：` + code + `，10分钟有效。如非本人操作，请忽略",
    "phone": "` + mobile + `",
    "sendtime": "201704101400",
    "report": "true",
    "extend": "555",
    "uid": "4321abc"
}`)

	fmt.Println(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ret := &SmsReturn{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	if ret.Code != "0" {
		return ret, errno.ErrCommon.WithMessage(ret.ErrorMsg)
	}
	return ret, nil
}
