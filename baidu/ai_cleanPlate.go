package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	aiCleanPlateUrl = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/emptyPlateRecognize"
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
	u := fmt.Sprintf("%s?access_token=%s&input_type=%s", aiCleanPlateUrl, token, "url")
	body, err := httptool.PostJson(u, map[string]string{
		"url": param.ImageUrl,
	})
	if err != nil {
		return nil, err
	}
	resp := CleanPlateResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
