package hellobike

import (
	"bytes"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	path       = "https://openapi.hellobike.com/bike/activity"
	AppId      = "20230302145050102"
	AppKey     = "d9244321dc3246caa54a29e7c156dd0c"
	activityId = "H3979885952972083867"
)

type Client struct {
	//https://app.trtpazyz.com
	Domain    string
	AppId     string
	Version   string
	htpClient http.Client
}

// TicketAllot 发放电子票
func (c *Client) SendCoupon(param SendCouponParam) (resp *BaseResponse, bizId string, err error) {
	//c.Domain+path
	return c.request("https://openapi.hellobike.com/bike/activity", param)
}

func (c *Client) request(url string, v interface{}) (resp *BaseResponse, bizId string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}
	s, err := c.sign(data)
	if err != nil {
		return nil, "", err
	}
	fmt.Println(string(data))
	fmt.Println(s)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, "", err
	}
	bizId = time.Now().Format("20060102150405") + c.rand()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appid", c.AppId)
	req.Header.Add("sequence", bizId)
	req.Header.Add("version", c.Version)
	req.Header.Add("signature", s)

	htpRes, err := c.htpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer htpRes.Body.Close()

	resp = &BaseResponse{}
	err = json.NewDecoder(htpRes.Body).Decode(resp)
	if err != nil {
		return nil, "", err
	}
	return resp, bizId, nil
}

func (c *Client) sign(v []byte) (string, error) {

	return "", nil
}

func (c *Client) rand() string {
	s := ""
	for i := 0; i < 10; i++ {
		s += strconv.Itoa(mrand.Intn(10))
	}
	return s
}
