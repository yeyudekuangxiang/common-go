package wxapp

import (
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"io/ioutil"
	"mio/pkg/wxapp/httputil"
)

func NewClient(c *weapp.Client) *Client {
	return &Client{Client: c}
}

type Client struct {
	*weapp.Client
}

// GetUnlimitedQRCodeResponse 获取没有数量限制的小程序码
func (c Client) GetUnlimitedQRCodeResponse(param *weapp.UnlimitedQRCode) (*QRCodeResponse, error) {
	resp, cerr, err := c.GetUnlimitedQRCode(param)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return &QRCodeResponse{
		Response: Response{
			ErrCode: cerr.ErrCode,
			ErrMsg:  cerr.ErrMSG,
		},
		ContentType: resp.Header.Get("Content-Type"),
		Buffer:      body,
	}, nil
}

// GetWxaCodeResponse 获取有数量限制的小程序码
func (c Client) GetWxaCodeResponse(code *weapp.QRCode) (*QRCodeResponse, error) {
	resp, cerr, err := c.GetQRCode(code)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return &QRCodeResponse{
		Response: Response{
			ErrCode: cerr.ErrCode,
			ErrMsg:  cerr.ErrMSG,
		},
		ContentType: resp.Header.Get("Content-Type"),
		Buffer:      body,
	}, nil
}

// CreateWxaQrcodeResponse 获取有数量限制的小程序二维码
func (c Client) CreateWxaQrcodeResponse(creator *weapp.QRCodeCreator) (*QRCodeResponse, error) {
	resp, cerr, err := c.CreateQRCode(creator)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return &QRCodeResponse{
		Response: Response{
			ErrCode: cerr.ErrCode,
			ErrMsg:  cerr.ErrMSG,
		},
		ContentType: resp.Header.Get("Content-Type"),
		Buffer:      body,
	}, nil
}

var accessToken = ""

// GetUserRiskRank 根据提交的用户信息获取用户的安全等级 recursiveCount记录token1失效时递归次数 最多重试三次
func (c Client) GetUserRiskRank(param UserRiskRankParam, recursiveCount ...int) (*UserRiskRankResponse, error) {
	if accessToken == "" {
		token, err := c.AccessToken()
		if err != nil {
			return nil, err
		}
		accessToken = token
	}

	u := fmt.Sprintf("https://api.weixin.qq.com/wxa/getuserriskrank?access_token=%s", accessToken)
	body, err := httputil.PostJson(u, param)
	if err != nil {
		return nil, err
	}

	resp := UserRiskRankResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode == 40001 || resp.ErrCode == 41001 || resp.ErrCode == 42001 {

		count := 0
		if len(recursiveCount) > 0 {
			count = recursiveCount[0]
		}

		if count >= 3 {
			return &resp, nil
		}
		count++

		token, err := c.AccessToken()
		if err != nil {
			return nil, err
		}
		accessToken = token
		return c.GetUserRiskRank(param, count)
	}
	return &resp, nil
}
