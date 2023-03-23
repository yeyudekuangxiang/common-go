package baidu

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"strconv"
)

const (
	imageReviewUrl = "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_defined"
	textReviewUrl  = "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined"
)

type ReviewClient struct {
	AccessToken IAccessToken
}

func NewReviewClient(accessToken IAccessToken) *ReviewClient {
	return &ReviewClient{AccessToken: accessToken}
}

// ImageReviewParam 入参
type ImageReviewParam struct {
	//图片base64
	Image string `json:"image,omitempty"`
	//图片链接
	ImgUrl string `json:"imgUrl,omitempty"`
	//图片类型0:静态图片（PNG、JPG、JPEG、BMP、GIF（仅对首帧进行审核）、Webp、TIFF），1:GIF动态图片
	ImgType uint64 `json:"imgType"`
}

type ImageReviewData struct {
	Type           int         `json:"type"`
	SubType        int         `json:"subType"`
	Conclusion     string      `json:"conclusion,omitempty"`
	ConclusionType int         `json:"conclusionType,omitempty"`
	Msg            string      `json:"msg"`
	Probability    float64     `json:"probability,omitempty"`
	Codes          []string    `json:"codes,omitempty"`
	DatasetName    string      `json:"datasetName,omitempty"`
	Completeness   float64     `json:"completeness,omitempty"`
	Hits           interface{} `json:"hits,omitempty"`
	Conclution     string      `json:"conclution,omitempty"`
	ConclutionType int         `json:"conclutionType,omitempty"`
}

// ImageReviewResp 出参
type ImageReviewResp struct {
	CommonRespCode
	LogId          int64             `json:"log_id"`
	Conclusion     string            `json:"conclusion"`
	ConclusionType int               `json:"conclusionType"`
	Data           []ImageReviewData `json:"data"`
}

// ImageReview  图片审核
func (l *ReviewClient) ImageReview(param ImageReviewParam) (*ImageReviewResp, error) {
	token, err := l.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s?access_token=%s", imageReviewUrl, token)
	body, err := httptool.PostMapFrom(u, map[string]string{
		"image":    param.Image,
		"imageUrl": param.ImgUrl,
		"imgType":  strconv.FormatUint(param.ImgType, 10),
	})
	if err != nil {
		return nil, err
	}
	resp := ImageReviewResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type TextReviewParam struct {
	Text string `json:"text"`
}
type TextReviewData struct {
	Type           int         `json:"type"`
	SubType        int         `json:"subType"`
	Conclusion     string      `json:"conclusion,omitempty"`
	ConclusionType int         `json:"conclusionType,omitempty"`
	Msg            string      `json:"msg"`
	Probability    float64     `json:"probability,omitempty"`
	Codes          []string    `json:"codes,omitempty"`
	DatasetName    string      `json:"datasetName,omitempty"`
	Completeness   float64     `json:"completeness,omitempty"`
	Hits           interface{} `json:"hits,omitempty"`
	Conclution     string      `json:"conclution,omitempty"`
	ConclutionType int         `json:"conclutionType,omitempty"`
}

// TextReviewResp 出参
type TextReviewResp struct {
	CommonRespCode
	LogId          int64             `json:"log_id"`
	Conclusion     string            `json:"conclusion"`
	ConclusionType int               `json:"conclusionType"`
	Data           []ImageReviewData `json:"data"`
}

// TextReview 文字审核
func (l *ReviewClient) TextReview(param TextReviewParam) (*TextReviewResp, error) {
	b, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}

	token, err := l.AccessToken.GetToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s?access_token=%s", textReviewUrl, token)
	body, err := httptool.PostJson(u, b)
	if err != nil {
		return nil, err
	}

	resp := TextReviewResp{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
