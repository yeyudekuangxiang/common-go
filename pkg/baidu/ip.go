package baidu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	if ip == "" {
		return nil, errors.New("ip地址为空")
	}
	url := "https://api.map.baidu.com/location/ip?ak=32f38c9491f2da9eb61106aaab1e9739&ip=" + ip
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Cookie", "")
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
	var o CityResult
	err = json.Unmarshal([]byte(string(body)), &o)
	if err != err {
		return nil, err
	}
	return &o, nil
}
