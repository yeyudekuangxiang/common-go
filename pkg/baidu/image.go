package baidu

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/util/httputil"
)

const (
	webimageUrl      = "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage"
	accurateBasicUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
)

type ImageClient struct {
	AccessToken *AccessToken
}
type WebImageParam struct {
	ImageUrl string
}
type ImageRespCode struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (result ImageRespCode) IsSuccess() bool {
	return result.ErrorCode == 0
}

type WebImageResult struct {
	ImageRespCode
	LogId          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

// WebImage 网络图片文字识别
func (c *ImageClient) WebImage(param WebImageParam) (*WebImageResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s", webimageUrl, token)
	body, err := httputil.PostMapFrom(u, map[string]string{
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
	ImageRespCode
	LogId          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

// AccurateBasic 通用文字识别（高精度版)
func (c *ImageClient) AccurateBasic(param AccurateBasicParam) (*AccurateBasicResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s", accurateBasicUrl, token)

	body, err := httputil.PostMapFrom(u, map[string]string{
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
