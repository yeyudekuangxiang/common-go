package Fmy_test

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/util/encrypt"
	"testing"
)

type data struct {
	PromOrderId int    `json:"prom_order_id,omitempty"`
	Phone       string `json:"phone,omitempty"`
	OutOrderNo  string `json:"out_order_no,omitempty"`
}

func TestFmy_GetSign(t *testing.T) {
	d := data{
		PromOrderId: 59,
		Phone:       "13688888888",
		OutOrderNo:  "2021111015441133138101",
	}
	marshal, _ := json.Marshal(d)
	fmt.Println(string(marshal))
	rand1 := "1111"
	rand2 := "2222"
	platformKey := "6ED925249F5E266892EE74243118D9D5"
	appSecret := "00w6upJInLtuwYAba5XbeKtxAucAMuX0"
	verifyData := rand1 + platformKey + string(marshal) + appSecret + rand2
	md5Str := encrypt.Md5(verifyData)
	sign := rand1 + md5Str[7:21] + rand2
	fmt.Println(sign)
}
