package message

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendYZM(mobile string, code string) {

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
