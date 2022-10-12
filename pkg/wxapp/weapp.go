package wxapp

import (
	"encoding/json"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"mio/pkg/wxapp/httputil"
	"strconv"
	"time"
)

type noCache struct {
}

func (noCache) Set(key string, val interface{}, timeout time.Duration) {
}

func (noCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func NewClient(appId, appSecret string, tokenCenter AccessTokenCenter) *Client {
	return &Client{Client: weapp.NewClient(appId, appSecret, weapp.WithAccessTokenSetter(func() (token string, expireIn uint) {
		token, expireAt, err := tokenCenter.AccessToken()

		if err != nil {
			log.Println("获取token失败", appId, appSecret, err)
			return "", 0
		}

		return token, uint(expireAt.Sub(time.Now()).Seconds())
	}), weapp.WithCache(noCache{})), tokenCenter: tokenCenter}
}

type Client struct {
	appId     string
	appSecret string
	*weapp.Client
	tokenCenter AccessTokenCenter
	Logger      *zap.SugaredLogger
}

// GetUnlimitedQRCodeResponse 获取没有数量限制的小程序码
func (c *Client) GetUnlimitedQRCodeResponse(param *weapp.UnlimitedQRCode) (*QRCodeResponse, error) {

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
func (c *Client) GetWxaCodeResponse(code *weapp.QRCode) (*QRCodeResponse, error) {
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
func (c *Client) CreateWxaQrcodeResponse(creator *weapp.QRCodeCreator) (*QRCodeResponse, error) {
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

// GetUserRiskRank 根据提交的用户信息获取用户的安全等级 recursiveCount记录token1失效时递归次数 最多重试三次
func (c *Client) GetUserRiskRank(param UserRiskRankParam, recursiveCount ...int) (*UserRiskRankResponse, error) {
	accessToken, err := c.AccessToken()
	if err != nil {
		return nil, err
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
	return &resp, nil
}

func (c *Client) IsExpireAccessToken(code int) (bool, error) {
	if code == 0 {
		return false, nil
	}
	return c.tokenCenter.IsExpired(strconv.Itoa(code))
}

func (c *Client) AccessToken() (string, error) {
	token, _, err := c.tokenCenter.AccessToken()
	if err != nil {
		return "", err
	}
	return token, nil
}
func (c *Client) AccessTokenFull() (string, time.Duration, error) {
	token, expireAt, err := c.tokenCenter.AccessToken()
	if err != nil {
		return "", 0, err
	}
	return token, expireAt.Sub(time.Now()), nil
}
func (c *Client) RefreshToken(oldToken string) (string, time.Duration, error) {
	token, expireAt, err := c.tokenCenter.RefreshToken(oldToken)
	if err != nil {
		return "", 0, err
	}
	return token, expireAt.Sub(time.Now()), nil
}

// AutoTryAccessToken 自动重试 获取token异常时会返回err
func (c *Client) AutoTryAccessToken(f func(accessToken string) (try bool, err error), maxTryCount int) error {
	accessToken, err := c.AccessToken()
	if err != nil {
		return err
	}

	try, err := f(accessToken)
	if err != nil {
		return err
	}

	if !try {
		return nil
	}

	if maxTryCount > 0 {
		_, _, err := c.RefreshToken(accessToken)
		if err != nil {
			return err
		}
		maxTryCount--
		return c.AutoTryAccessToken(f, maxTryCount)
	}
	return nil
}
