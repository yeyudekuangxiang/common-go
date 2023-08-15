package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	aiCleanPlate = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/emptyPlateRecognize" //光盘
	aiRecycle    = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/recycle"             //旧瓶
	aiWRS        = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/classification/lvmio_wrs"           //垃圾分类
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
	ImageUrl string //url
	Image    string //base64
}

type AiImageResult struct {
	CommonRespCode
	LogId   int64    `json:"log_id"`
	Results []Result `json:"results"`
}

type Result struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func (c *AiImageClient) CleanPlate(param AiImageParam) (*AiImageResult, error) {
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

	resp := AiImageResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *AiImageClient) GarbageSorting(param AiImageParam) (*AiImageResult, error) {
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
	resp := AiImageResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
