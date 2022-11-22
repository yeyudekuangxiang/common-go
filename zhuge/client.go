package zhuge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	types2 "gitlab.miotech.com/miotech-application/backend/common-go/zhuge/types"
	"net/http"
	"time"
)

const EventUrl = "https://u.zhugeapi.com/open/v2/event_statis_srv/upload_event"

// Client 参考文档https://docs.zhugeio.com/dev/server2.html
type Client struct {
	appKey     string
	secretKey  string
	httpClient http.Client
	//1开启调试 0关闭调试
	debug int
}

func NewClient(appKey, secretKey string, debug int) *Client {
	return &Client{appKey: appKey, secretKey: secretKey, debug: debug, httpClient: http.Client{Timeout: 10 * time.Second}}
}

func (client Client) sign() string {
	signStr := client.appKey + ":" + client.secretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(signStr))
}
func (client Client) Track(event types2.Event, others map[string]interface{}) error {
	var pr map[string]interface{}
	if event.Pr != nil {
		pr = event.Pr.ToMap()
	} else {
		pr = make(map[string]interface{})
	}

	for k, v := range others {
		pr[k] = v
	}
	event.Pr = types2.MapPr(pr)

	fullEvent := types2.EventWithAk{
		Event: event,
		Ak:    client.appKey,
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
	req.Header.Set("Authorization", client.sign())

	resp, err := client.httpClient.Do(req)
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
func (client Client) TrackSimple(eventName, openId string, attr map[string]interface{}) error {
	err := client.Track(types2.Event{
		Dt:    "evt",
		Pl:    "js",
		Debug: client.debug,
		Pr: types2.EventJs{
			Ct:   time.Now().UnixMilli(),
			Eid:  eventName,
			Cuid: openId,
			Sid:  time.Now().UnixMilli(),
		},
	}, attr)

	return err
}

type Resp struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
	WarnDid       string `json:"warn_did"`
}
