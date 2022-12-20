package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"mio/internal/pkg/util/httputil"
)

const (
	imageReviewUrl = "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_defined"
	textReviewUrl  = "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined"
)

type ReviewClient struct {
	AccessToken *AccessToken
}

// ImageReviewParam 入参
type ImageReviewParam struct {
	Image    string `json:"image"`
	ImageUrl string `json:"imageUrl"`
}

// ReviewResp 出参
type ReviewResp struct {
	LogId          int64     `json:"log_id"`
	ErrorCode      int64     `json:"error_code,omitempty"`
	ErrorMsg       string    `json:"error_msg,omitempty"`
	Conclusion     string    `json:"conclusion"`
	ConclusionType int64     `json:"conclusionType"`
	Data           []DataRes `json:"data"`
}

type DataRes struct {
	TP             int    `json:"type"`
	SubType        int    `json:"subType"`
	Conclusion     string `json:"conclusion"`
	ConclusionType int64  `json:"conclusionType"`
	Msg            string `json:"msg"`
}

// ImageReview  图片审核
func (l *ReviewClient) ImageReview(param ImageReviewParam) (ReviewResp, error) {
	resp := ReviewResp{}
	token, _ := l.AccessToken.GetToken()

	u := fmt.Sprintf("%s?access_token=%s", imageReviewUrl, token)

	body, err := httputil.PostMapFrom(u, map[string]string{
		"imgUrl": param.ImageUrl,
		"image":  param.Image,
	})

	if err != nil {
		logx.Errorf("image review err: %s", err.Error())
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
