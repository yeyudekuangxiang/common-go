package ytx

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"strconv"
	"time"
)

//openid:  CpziorTGUL02NrrBqsbbhsAN0Ve4ZMSpPEmgBPAGZOY=
//secret:   a123456
//appid: cc5dec82209c45888620eabec3a29b50
//poolCode: RP202110251300002

const (
	grantV2 = "/markting_redenvelopegateway/redenvelope/grantV2"
)

type Client struct {
	Domain   string
	Secret   string
	PoolCode string
	AppId    string
}

func (c *Client) SendCoupon(orderId, thirdUserId string, amount float64, channelKey string) (*GrantV2Response, error) {
	grantV2Request := GrantV2Request{
		AppId:     c.AppId,
		AppSecret: c.getAppSecret(),
		Ts:        strconv.FormatInt(time.Now().Unix(), 10),
		ReqData: GrantV2ReqData{
			OrderNo:  orderId,
			PoolCode: c.PoolCode,
			Amount:   amount,
			OpenId:   thirdUserId,
			Remark:   fmt.Sprintf("%s%f%s", channelKey, amount, "元红包"),
		},
	}

	url := c.Domain + grantV2
	body, err := httptool.PostJson(url, grantV2Request)
	if err != nil {
		return nil, err
	}

	response := GrantV2Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) getAppSecret() string {
	return encrypttool.Md5(c.AppId + encrypttool.Md5(c.Secret) + strconv.FormatInt(time.Now().Unix(), 10))
}
