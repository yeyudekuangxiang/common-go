package auth

import (
	"encoding/json"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/idtool"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	URL = "https://authn.uat.sheca.com/dyrz/openApi/auth/mobile"
)

type Client struct {
	UniTrustAppId string
	Token         string
}

func (c *Client) SendAuth(req UserIdentityVerificationReq) (*UserIdentityVerificationResp, error) {
	param, err := json.Marshal(map[string]string{
		"name":          req.Name,
		"idNo":          req.IdentityCard,
		"mobile":        req.Phone,
		"transactionId": idtool.UUID(),
	})
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(string(param))
	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", URL, payload)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Authorization", c.Token)
	httpReq.Header.Add("Content-Type", "application/json; charset=UTF-8")
	httpReq.Header.Add("UniTrust-AppId", c.UniTrustAppId)
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp2 := UserIdentityVerificationResp{}
	err = json.Unmarshal(body, &resp2)
	if err != nil {
		return nil, err
	}
	return &resp2, nil
}