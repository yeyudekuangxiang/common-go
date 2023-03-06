package miosass

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/sorttool"
	"log"
	"net/url"
	"strconv"
)

const (
	RedeemCertificateURL = "/api/mp2c/spu/public-welfare/certificate/redeem"
	CertificateCountURL  = "/api/mp2c/spu/public-welfare/certificate/mio/individual/record/count"
	AutoLoginURL         = "/api/mp2c/login/mio/login"
)

type Client struct {
	Domain    string
	AppKey    string
	AccessKey string
}

func (c *Client) CertificateCount(param CertificateCountParam) (*CertificateCountResp, error) {
	u := fmt.Sprintf("%s%s", c.Domain, CertificateCountURL)
	param.AppKey = c.AppKey
	param.Sign = c.sign(param.signParams())

	uv := url.Values{}
	for k, v := range param.signParams() {
		uv.Set(k, v)
	}
	uv.Set("sign", param.Sign)

	log.Printf("获取证书数量 %v", uv)
	respBody, err := httptool.Get(u + "?" + uv.Encode())
	if err != nil {
		return nil, err
	}
	resp := CertificateCountResp{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) ExchangeCertificate(param ExchangeCertificateParam) (*ExchangeCertificateResp, error) {
	u := fmt.Sprintf("%s%s", c.Domain, RedeemCertificateURL)
	param.AppKey = c.AppKey
	param.Sign = c.sign(param.signParams())
	body, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	log.Printf("兑换证书 %v", string(body))
	respBody, err := httptool.PostJsonBytes(u, body)
	if err != nil {
		return nil, err
	}
	resp := ExchangeCertificateResp{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (c *Client) AutoLogin(param AutoLoginParam) string {
	param.addAppKey(c.AppKey)
	signStr := c.sign(param.signParams())
	signPrams := param.signParams()
	signPrams["sign"] = signStr

	p := url.Values{}

	sorttool.Map(signPrams, func(key interface{}) {
		p.Add(key.(string), signPrams[key.(string)])
	})
	return fmt.Sprintf("%s%s?%s", c.Domain, AutoLoginURL, p.Encode())
}
func (c *Client) sign(signMap map[string]string) string {
	signStr := ""
	sorttool.Map(signMap, func(key interface{}) {
		keyStr := key.(string)
		signStr += keyStr + "=" + signMap[keyStr]
	})
	log.Println("sass验签", c.AccessKey+signStr)
	return encrypttool.Md5(c.AccessKey + signStr)
}
func (c *Client) Sign(param map[string]string) string {
	return c.sign(param)
}

type CertificateCountResp struct {
	ErrResp
	Data struct {
		Count int64 `json:"count"`
	}
}
type CertificateCountParam struct {
	sign
	UserId string `json:"userId"`
}

func (c CertificateCountParam) signParams() map[string]string {
	return map[string]string{
		"appKey": c.AppKey,
		"userId": c.UserId,
	}
}

type ExchangeCertificateParam struct {
	sign
	SkuId   int64  `json:"skuId"`
	UserId  string `json:"userId"`
	HeadImg string `json:"headImg"`
	Name    string `json:"name"`
}

func (c *ExchangeCertificateParam) signParams() map[string]string {
	return map[string]string{
		"skuId":   strconv.FormatInt(c.SkuId, 10),
		"userId":  c.UserId,
		"headImg": c.HeadImg,
		"name":    c.Name,
		"appKey":  c.AppKey,
	}
}

type ExchangeCertificateResp struct {
	ErrResp
	Data struct {
		BizId          string `json:"bizId"`
		CertificateUrl string `json:"certificateUrl"`
	}
}

type ErrResp struct {
	Code       int    `json:"code"`
	DetailCode string `json:"detailCode"`
	Message    string `json:"message"`
}

func (r ErrResp) IsSuccess() bool {
	return r.Code == 200
}

type AutoLoginParam struct {
	sign
	Name    string `json:"name"`
	Credits int64  `json:"credits"`
	HeadImg string `json:"headImg"`
	Uid     string `json:"uid"`
	//毫秒时间戳
	Timestamp int64 `json:"timestamp"`
}

func (c *AutoLoginParam) signParams() map[string]string {
	return map[string]string{
		"name":      c.Name,
		"credits":   strconv.FormatInt(c.Credits, 10),
		"headImg":   c.HeadImg,
		"uid":       c.Uid,
		"timestamp": strconv.FormatInt(c.Timestamp, 10),
		"appKey":    c.AppKey,
	}
}

type sign struct {
	Sign   string `json:"sign"`
	AppKey string `json:"appKey"`
}

func (s *sign) addSign(sign string) {
	s.Sign = sign
}
func (s *sign) addAppKey(appKey string) {
	s.AppKey = appKey
}
