package wxapp

import (
	"github.com/medivhzhan/weapp/v3"
	"github.com/pkg/errors"
	"io/ioutil"
)

func NewClient(c *weapp.Client) *Client {
	return &Client{c: c}
}

type Client struct {
	c *weapp.Client
}

func (c Client) GetUnlimitedQRCode(param *weapp.UnlimitedQRCode) (*UnlimitedQRCodeResponse, error) {
	resp, cerr, err := c.c.GetUnlimitedQRCode(param)
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
