package gaode

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/baidu"
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
type MapWithBaiduClient struct {
	GaoDeAk string
	BaiDuAk string
}

//优先使用高德,高德查不到，再使用百度

func NewMapWithBaiduClient(gaoDeAk string, baiduAk string) *MapWithBaiduClient {
	return &MapWithBaiduClient{GaoDeAk: gaoDeAk, BaiDuAk: baiduAk}
}

func (c *MapWithBaiduClient) LocationIp(ip string) (*LocationIpResultV2, error) {
	locationIp, err := NewMapClient(c.GaoDeAk).LocationIp(ip)
	if err != nil {
		return nil, err
	}
	if locationIp.Status == "1" && locationIp.Province != "" {
		return &LocationIpResultV2{
			StatusRespCode: StatusRespCode{
				Status: locationIp.Status,
			},
			Info:      locationIp.Info,
			Infocode:  locationIp.Infocode,
			Province:  locationIp.Province,
			City:      locationIp.City,
			Adcode:    locationIp.Adcode,
			Rectangle: locationIp.Rectangle,
		}, nil
	}
	result, err := baidu.NewMapClient(c.BaiDuAk).LocationIp(ip)
	if err != nil {
		return nil, err
	}
	if result.IsSuccess() {
		return &LocationIpResultV2{
			StatusRespCode: StatusRespCode{
				Status: "1",
			},
			Info:      "ok",
			Infocode:  "",
			Province:  result.Content.AddressDetail.Province,
			City:      result.Content.AddressDetail.City,
			Adcode:    result.Content.AddressDetail.Adcode,
			Rectangle: "",
		}, nil
	}
	return nil, errors.New("ip定位失败")
}

type StatusRespCode struct {
	Status string `json:"status"`
}

func (result StatusRespCode) IsSuccess() bool {
	return result.Status == "1"
}

type LocationIpResultV2 struct {
	StatusRespCode
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}
