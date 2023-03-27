package jhx

import (
	"encoding/json"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ticket = "/busticket/ticket_create"
)

type Client struct {
	Domain  string
	AppId   string
	Version string
}

//发放券码
func (c *Client) SendCoupon(orderId, phone string) (*CommonResponse, error) {
	commonParams := make(map[string]interface{}, 0)
	commonParams["version"] = c.Version
	commonParams["appid"] = c.AppId
	commonParams["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	commonParams["nonce"] = strconv.Itoa(rand.Int())
	commonParams["tradeno"] = orderId
	commonParams["mobile"] = phone
	commonParams["sign"] = strings.ToUpper(c.getSign(commonParams, "", "&"))

	url := c.Domain + ticket
	body, err := httptool.PostJson(url, commonParams)
	if err != nil {
		return nil, err
	}

	response := CommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) getSign(params map[string]interface{}, key string, joiner string) string {
	if joiner == "" {
		joiner = ";"
	}
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + c.interfaceToString(params[v]) + joiner
	}
	if joiner != ";" {
		signStr = strings.TrimRight(signStr, joiner)
	}
	//验证签名
	return encrypttool.Md5(key + signStr)
}

func (c *Client) interfaceToString(data interface{}) string {
	var key string
	switch data.(type) {
	case string:
		key = data.(string)
	case int:
		key = strconv.Itoa(data.(int))
	case int64:
		it := data.(int64)
		key = strconv.FormatInt(it, 10)
	case float64:
		it := data.(float64)
		key = strconv.FormatFloat(it, 'f', -1, 64)
	}
	return key
}
