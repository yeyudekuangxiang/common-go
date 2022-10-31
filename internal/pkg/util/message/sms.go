package message

import (
	"fmt"
	"io/ioutil"
	"mio/config"
	"net/http"
	"strings"
)

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
