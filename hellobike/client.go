package hellobike

import (
	"bytes"
	"encoding/json"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	mrand "math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	path       = "https://openapi.hellobike.com/bike/activity"
	AppId      = "20230302145050102"
	AppKey     = "d9244321dc3246caa54a29e7c156dd0c"
	activityId = "H3979885952972083867"
)

type Client struct {
	AppId     string
	Version   string
	htpClient http.Client
	Action    string
	AppKey    string
	Domain    string
}

// SendCoupon 发放电子票

func (c *Client) SendCoupon(param SendCouponParam) (resp *BaseResponse, err error) {
	//c.Domain+path
	return c.request(c.Domain, param)
}

func (c *Client) request(url string, param SendCouponParam) (resp *BaseResponse, err error) {
	params := make(map[string]string, 0)
	bizContent, _ := json.Marshal(param.BizContent)
	params["version"] = param.Version
	params["action"] = param.Action
	params["app_id"] = param.AppId
	params["biz_content"] = string(bizContent)
	params["utc_timestamp"] = param.UtcTimestamp
	sign := c.GetSign(params, "&", c.AppKey)
	param.Sign = sign
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshal))
	if err != nil {
		return nil, err
	}
	htpRes, err := c.htpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer htpRes.Body.Close()

	resp = &BaseResponse{}
	err = json.NewDecoder(htpRes.Body).Decode(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) rand() string {
	s := ""
	for i := 0; i < 10; i++ {
		s += strconv.Itoa(mrand.Intn(10))
	}
	return s
}

// GetSign 签名
func (c *Client) GetSign(params map[string]string, joiner string, appKey string) string {
	if joiner == "" {
		joiner = "&"
	}
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + params[v] + joiner
	}
	if joiner != ";" {
		signStr = strings.TrimRight(signStr, joiner)
	}
	//验证签名
	return encrypttool.Md5(signStr + appKey)
}

func (c *Client) RefundCard(param RefundCardParam) (resp *RefundCardResponse, err error) {
	//c.Domain+path
	return c.RefundCardRequest(c.Domain, param)
}

func (c *Client) RefundCardRequest(url string, param RefundCardParam) (resp *RefundCardResponse, err error) {
	params := make(map[string]string, 0)
	bizContent, _ := json.Marshal(param.BizContent)
	params["version"] = param.Version
	params["action"] = param.Action
	params["appId"] = param.AppId
	params["bizContent"] = string(bizContent)
	params["timestamp"] = param.Timestamp
	sign := c.GetSign(params, "&", c.AppKey)
	param.Sign = sign
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshal))
	if err != nil {
		return nil, err
	}
	htpRes, err := c.htpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer htpRes.Body.Close()

	resp = &RefundCardResponse{}
	err = json.NewDecoder(htpRes.Body).Decode(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
