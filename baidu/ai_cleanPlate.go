package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	aiCleanPlateUrl = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/emptyPlateRecognize"
	aiRecycleUrl    = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/emptyPlateRecognize"
)

type AiCleanPlateClient struct {
	AccessToken IAccessToken
}

func NewAiCleanPlateClient(accessToken IAccessToken) *AiCleanPlateClient {
	return &AiCleanPlateClient{
		AccessToken: accessToken,
	}
}

type CleanPlateParam struct {
	ImageUrl string
	Image    string
}

type CleanPlateResult struct {
	CommonRespCode
	LogId   int64 `json:"log_id"`
	Results []struct {
		Name  string  `json:"name"`
		Score float64 `json:"score"`
	} `json:"results"`
}

func (c *AiCleanPlateClient) CleanPlate(param CleanPlateParam) (*CleanPlateResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s?access_token=%s", aiCleanPlate, token)
	p := make(map[string]string, 0)

	if param.Image != "" {
		p["image"] = param.Image
	}

	if param.ImageUrl != "" {
		u = fmt.Sprintf("%s?access_token=%s&input_type=%s", aiCleanPlate, token, "url")
		p["url"] = param.ImageUrl
	}

	body, err := httptool.PostJson(u, p)
	if err != nil {
		return nil, err
	}
	resp := CleanPlateResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *AiCleanPlateClient) Recycle(param AiImageParam) (*AiImageResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s&input_type=%s", aiRecycleUrl, token, "url")
	body, err := httptool.PostJson(u, map[string]string{
		"url": param.ImageUrl,
	})
	if err != nil {
		return nil, err
	}
	resp := AiImageResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
