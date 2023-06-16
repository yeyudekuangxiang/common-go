package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	aiCleanPlate = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/emptyPlateRecognize"
	aiRecycle    = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/recycle"
)

type AiImageClient struct {
	AccessToken IAccessToken
}

func NewAiImageClient(accessToken IAccessToken) *AiImageClient {
	return &AiImageClient{
		AccessToken: accessToken,
	}
}

type AiImageParam struct {
	ImageUrl string
	Addr     string
}

type AiImageResult struct {
	CommonRespCode
	LogId   int64 `json:"log_id"`
	Results []struct {
		Name  string  `json:"name"`
		Score float64 `json:"score"`
	} `json:"results"`
}

func (c *AiImageClient) CleanPlate(param AiImageParam) (*AiImageResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s&input_type=%s", aiCleanPlate, token, "url")
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

func (c *AiImageClient) Recycle(param AiImageParam) (*AiImageResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s&input_type=%s", aiRecycle, token, "url")
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
