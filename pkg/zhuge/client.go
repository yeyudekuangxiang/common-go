package zhuge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mio/pkg/zhuge/types"
	"net/http"
	"time"
)

const EventUrl = "https://u.zhugeapi.com/open/v2/event_statis_srv/upload_event"

// Client 参考文档https://docs.zhugeio.com/dev/server2.html
type Client struct {
	AppKey    string
	SecretKey string
	htpClient http.Client
}

func NewClient(appKey, secretKey string) *Client {
	return &Client{AppKey: appKey, SecretKey: secretKey, htpClient: http.Client{Timeout: 10 * time.Second}}
}

func (client Client) Sign() string {
	signStr := client.AppKey + ":" + client.SecretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(signStr))
}
func (client Client) Track(event types.Event, others map[string]interface{}) error {
	var pr map[string]interface{}
	if event.Pr != nil {
		pr = event.Pr.ToMap()
	} else {
		pr = make(map[string]interface{})
	}

	for k, v := range others {
		pr[k] = v
	}
	event.Pr = types.MapPr(pr)

	fullEvent := types.EventWithAk{
		Event: event,
		Ak:    client.AppKey,
	}

	body, err := json.Marshal(fullEvent)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", EventUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", client.Sign())

	resp, err := client.htpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	zhugeResp := Resp{}

	err = json.NewDecoder(resp.Body).Decode(&zhugeResp)
	if err != nil {
		return err
	}

	if zhugeResp.ReturnCode != 0 {
		return errors.New(fmt.Sprintf("code:%d message:%s warn:%s", zhugeResp.ReturnCode, zhugeResp.ReturnMessage, zhugeResp.WarnDid))
	}
	return nil
}

type Resp struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	WarnDid       string `json:"warn_did"`
}
