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

func (c Client) GetUnlimitedQRCodeResponse(param *weapp.UnlimitedQRCode) (*UnlimitedQRCodeResponse, error) {
	resp, cerr, err := c.GetUnlimitedQRCode(param)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	return &UnlimitedQRCodeResponse{
		Response: Response{
			ErrCode: cerr.ErrCode,
			ErrMsg:  cerr.ErrMSG,
		},
		ContentType: resp.Header.Get("Content-Type"),
		Buffer:      body,
	}, nil
}

var accessToken = ""

// GetUserRiskRank 根据提交的用户信息获取用户的安全等级
func (c Client) GetUserRiskRank(param UserRiskRankParam) (*UserRiskRankResponse, error) {
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
	fmt.Println("request", string(body))
	resp := UserRiskRankResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode == 40001 {
		token, err := c.AccessToken()
		if err != nil {
			return nil, err
		}
		accessToken = token
		return c.GetUserRiskRank(param)
	}
	return &resp, nil
}
