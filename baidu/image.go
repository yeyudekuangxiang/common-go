package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const (
	webimageUrl      = "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage"
	accurateBasicUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
	classifyBasicUrl = "https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general"
)

type ImageClient struct {
	AccessToken *AccessToken
}

func NewImageClient(accessToken *AccessToken) *ImageClient {
	return &ImageClient{
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
func (c *ImageClient) WebImage(param WebImageParam) (*WebImageResult, error) {
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
func (c *ImageClient) AccurateBasic(param AccurateBasicParam) (*AccurateBasicResult, error) {
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

// AdvancedGeneralParam 通用物体和场景识别入参
type AdvancedGeneralParam struct {
	Image    string `json:"image"`     //base64格式
	Url      string `json:"url"`       //url格式
	BaikeNum int    `json:"baike_num"` //用于控制返回结果是否带有百科信息 若不输入此参数，则默认不返回百科结果；若输入此参数，会根据输入的整数返回相应个数的百科信息
}

// AdvancedGeneralResult 通用物体和场景识别出参
type AdvancedGeneralResult struct {
	LogId     uint64 `json:"log_id"`
	ResultNum uint32 `json:"result_num"`
	Result    []struct {
		Keyword   string  `json:"keyword"`
		Score     float64 `json:"score"`
		Root      string  `json:"root"`
		BaikeInfo []struct {
			BaikeUrl    string `json:"baike_url"`
			ImageUrl    string `json:"image_url"`
			Description string `json:"description"`
		}
	}
}

// AdvancedGeneral 通用物体和场景识别
func (c *ImageClient) AdvancedGeneral(param AdvancedGeneralParam) (*AdvancedGeneralResult, error) {
	token, err := c.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s", classifyBasicUrl, token)

	body, err := httptool.PostMapFrom(u, map[string]string{
		"url": param.Url,
	})
	if err != nil {
		return nil, err
	}
	resp := AdvancedGeneralResult{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
