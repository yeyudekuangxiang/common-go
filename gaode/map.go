package gaode

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

type MapClient struct {
	ak string
}

func NewMapClient(ak string) *MapClient {
	return &MapClient{ak: ak}
}

func (c *MapClient) LocationIp(ip string) (*LocationIpResult, error) {
	url := fmt.Sprintf("https://restapi.amap.com/v3/ip?key=%s&ip=%s", c.ak, ip)
	body, err := httptool.Get(url)
	if err != nil {
		return nil, err
	}

	var o LocationIpResult
	err = json.Unmarshal(body, &o)
	if err != err {
		return nil, err
	}
	return &o, nil
}

type LocationIpResult struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

type T struct {
	Info         string `json:"info"`
	Infocode     string `json:"infocode"`
	Status       string `json:"status"`
	SecCodeDebug string `json:"sec_code_debug"`
	Key          string `json:"key"`
	SecCode      string `json:"sec_code"`
}
