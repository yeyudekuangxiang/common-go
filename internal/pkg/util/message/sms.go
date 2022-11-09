package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chanxuehong/util/url"
	"io/ioutil"
	"mio/config"
	"mio/pkg/errno"
	"net/http"
	"strings"
	"unsafe"
)

type JsonPostSample struct {
}

func SendYZMV4(mobile string, code string) {

	url := "https://smssh1.253.com/msg/v1/send/json"
	method := "POST"

	payload := strings.NewReader(`{
    "account": "YZM7795025",
    "password": "P4tDNsDCXI5380", //需要加入K8S
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
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

var TemplateMap = map[string]string{
	"839944": "您发布的“{$var}”，买家已付款，等待您发货，戳 ezjf.cn/v-3dzyg{$var} 查看 退订回T",
	"839946": "您发布的“{$var}”，由于您超时未发货，交易取消，戳 ezjf.cn/v-3dzyg{$var}",
}

func SendYZMV5(templateId string, phone string, msg string) (string, error) {
	msg = phone + ", " + msg + ";"
	template, ok := TemplateMap[templateId]
	if !ok {
		return "", errno.ErrCommon.WithMessage("模版有误" + templateId)
	}
	url := "https://smssh1.253.com/msg/variable/json"
	method := "POST"
	payload := strings.NewReader(`{
    "account": "M4232956",
    "password": "8Xx53be5pXc568", //需要加入K8S
    "msg": "` + template + `",
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
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))

	return string(body), nil
}

func SendSmsV2(mobile string, sms string) error {
	params := make(map[string]interface{})
	//请登录zz.253.com获取API账号、密码以及短信发送的URL
	params["account"] = config.Config.Sms.Account   //创蓝API账号
	params["password"] = config.Config.Sms.Password //创蓝API密码
	params["phone"] = mobile                        //手机号码

	//设置您要发送的内容：其中“【】”中括号为运营商签名符号，多签名内容前置添加提交
	a := "验证码{686739}，30分钟有效。参与低碳任务，体验格调生活。如非本人操作请忽略。"
	params["msg"] = url.QueryEscape(a)

	params["report"] = "true"
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	reader := bytes.NewReader(bytesData)
	url := "https://smssh1.253.com/msg/v1/send/json" //短信发送URL
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
	return nil
}

func SendSms(mobile string, sms string) error {
	url := "https://smssh1.253.com/msg/v1/send/json"
	method := "POST"
	payload := strings.NewReader(`{
    "account": "` + config.Config.Sms.Account + `",
    "password": "` + config.Config.Sms.Password + `", //需要加入K8S
    "msg": "` + sms + `",
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
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
