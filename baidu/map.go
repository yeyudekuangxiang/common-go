package baidu

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
	url := fmt.Sprintf("https://api.map.baidu.com/location/ip?ak=%s&ip=%s", c.ak, ip)
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
	StatusRespCode
	Address string `json:"address"`
	Content struct {
		AddressDetail struct {
			Province     string `json:"province"`
			City         string `json:"city"`
			District     string `json:"district"`
			StreetNumber string `json:"street_number"`
			Adcode       string `json:"adcode"`
			Street       string `json:"street"`
			CityCode     int    `json:"city_code"`
		} `json:"address_detail"`
		Point struct {
			Y string `json:"y"`
			X string `json:"x"`
		} `json:"point"`
		Address string `json:"address"`
	} `json:"content"`
}
