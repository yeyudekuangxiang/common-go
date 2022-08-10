package baidu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CityResult struct {
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
	Status int `json:"status"`
}

//参考地址 https://www.cnblogs.com/finger-ghost/p/14303291.html
//根据ip获取城市

func IpToCity(ip string) (*CityResult, error) {
	url := "https://api.map.baidu.com/location/ip?ak=32f38c9491f2da9eb61106aaab1e9739"
	method := "GET"
	payload := strings.NewReader("ip=" + ip)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var o CityResult
	err = json.Unmarshal([]byte(string(body)), &o)
	fmt.Println("OCRPush ", payload, string(body))
	if err != err {
		return nil, err
	}
	return &o, nil
}
