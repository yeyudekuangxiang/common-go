package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	webimageUrl      = "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage"
	accurateBasicUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
)

type OCRClient struct {
	AccessToken *AccessToken
}

func NewOCRClient(accessToken *AccessToken) *OCRClient {
	return &OCRClient{
		AccessToken: accessToken,
	}
}

type WebImageParam struct {
	ImageUrl string
}

type WebImageResult struct {
	CommonRespCode
	LogId          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

// WebImage 网络图片文字识别
func (c *OCRClient) WebImage(param WebImageParam) (*WebImageResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s", webimageUrl, token)
	body, err := httptool.PostMapFrom(u, map[string]string{
		"url": param.ImageUrl,
	})
	if err != nil {
		return nil, err
	}
	resp := WebImageResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AccurateBasicParam struct {
	ImageUrl string
}

type AccurateBasicResult struct {
	CommonRespCode
	LogId          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

// AccurateBasic 通用文字识别（高精度版)
func (c *OCRClient) AccurateBasic(param AccurateBasicParam) (*AccurateBasicResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s", accurateBasicUrl, token)

	body, err := httptool.PostMapFrom(u, map[string]string{
		"url": param.ImageUrl,
	})
	if err != nil {
		return nil, err
	}
	resp := AccurateBasicResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
