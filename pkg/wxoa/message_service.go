package wxoa

import (
	"bytes"
	"encoding/json"
	"github.com/chanxuehong/wechat/mp/core"
	"io/ioutil"
	"net/http"
)

type MessageServer struct {
	AppId       string
	TokenServer core.AccessTokenServer
}

type MessageDataValueStruct struct {
	Value string `json:"value"`
}

type MessageDataStruct struct {
	Number1 MessageDataValueStruct `json:"number01"`
	Thing5  MessageDataValueStruct `json:"date01"`
	Time2   MessageDataValueStruct `json:"site01"`
	Number3 MessageDataValueStruct `json:"site02"`
}

type MessageStruct struct {
	Touser           string            `json:"touser"`
	TemplateId       string            `json:"template_id"`
	Page             string            `json:"page"`
	MiniprogramState string            `json:"miniprogram_state"`
	Lang             string            `json:"lang"`
	Data             MessageDataStruct `json:"data"`
}

func (srv *MessageServer) SendMessage() (*Ticket, error) {
	accessToken, err := srv.TokenServer.Token()
	println(accessToken)
	if err != nil {
		return nil, err
	}
	messDate := MessageStruct{}
	messDate.Touser = "oy_BA5FK1t3dEwrMZndhlUoI2-HY"
	messDate.TemplateId = "hQzUwkGMYqgNsKOad7RnIwwBpfkVfsuJvW6UqymwI8k"
	messDate.Page = "index"
	messDate.MiniprogramState = "developer"
	messDate.Lang = "zh_CN"
	messDate.Data.Number1.Value = "1"
	messDate.Data.Thing5.Value = "1"
	messDate.Data.Time2.Value = "1"
	messDate.Data.Number3.Value = "1"

	data, err := json.Marshal(messDate)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token="+accessToken, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, nil
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (srv *MessageServer) IIDB04E44A0E1DC11E5ADCEA4DB30FED8E1() {

}
