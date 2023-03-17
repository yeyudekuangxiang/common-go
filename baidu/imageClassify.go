package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
)

const classifyBasicUrl = "https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general"

type ClassifyClient struct {
	AccessToken *AccessToken
}

func NewClassifyClient(accessToken *AccessToken) *ClassifyClient {
	return &ClassifyClient{AccessToken: accessToken}
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
func (c *ClassifyClient) AdvancedGeneral(param AdvancedGeneralParam) (*AdvancedGeneralResult, error) {
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
