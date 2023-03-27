package ytx

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/idtool"
	"math/rand"
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
	ctx      context.Context
	Domain   string
	Secret   string
	PoolCode string
	AppId    string
}

func (c *Client) SendCoupon(thirdUserId string, amount float64) (*GrantV2Response, error) {
	rand.Seed(time.Now().UnixNano())
	grantV2Request := GrantV2Request{
		AppId:     c.AppId,
		AppSecret: c.getAppSecret(),
		Ts:        strconv.FormatInt(time.Now().Unix(), 10),
		ReqData: GrantV2ReqData{
			OrderNo:  "ytx" + idtool.UUID(),
			PoolCode: c.PoolCode,
			Amount:   amount,
			OpenId:   thirdUserId,
			Remark:   "lvmiao" + strconv.FormatFloat(amount, 'f', -1, 64) + "元红包",
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

	if response.SubCode != "0000" {
		return nil, errors.New(response.SubMessage)
	}
	return &response, nil
}

func (c *Client) getAppSecret() string {
	t := time.Now().Unix()
	return encrypttool.Md5(c.AppId + encrypttool.Md5(c.Secret) + strconv.FormatInt(t, 10))
}
