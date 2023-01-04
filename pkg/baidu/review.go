package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	"strings"
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
	//Image  string `json:"image,omitempty"`
	ImgUrl string `json:"imgUrl,omitempty"`
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
func (l *ReviewClient) ImageReview(param ImageReviewParam) error {
	token, err := l.AccessToken.GetToken()

	if err != nil {
		return err
	}

	u := fmt.Sprintf("%s?access_token=%s", imageReviewUrl, token)

	imageUrls := strings.Split(strings.Trim(param.ImgUrl, ","), ",")
	for _, url := range imageUrls {
		m := map[string]string{
			"imgUrl": url,
		}
		body, err := httputil.PostMapFrom(u, m)

		if err != nil {
			return errno.ErrCheckErr.WithMessage(fmt.Sprintf("系统错误: %s", err.Error()))
		}
		resp := &ReviewResp{}
		if err = json.Unmarshal(body, resp); err != nil {
			return errno.ErrCheckErr.WithMessage(fmt.Sprintf("系统错误: %s", err.Error()))
		}

		if resp.ErrorMsg != "" {
			app.Logger.Infof("review err : image_review param is %s, resp is %v", param, resp)
			return errno.ErrCheckErr.WithMessage(fmt.Sprintf("系统错误: %s", resp.ErrorMsg))
		}

		if resp.ConclusionType == 4 {
			return errno.ErrCheckErr.WithMessage("审核失败")
		}

		if resp.ConclusionType != 1 {
			return errno.ErrCheckErr.WithMessage(resp.Data[0].Msg)
		}
	}

	return nil
}

type TextReviewParam struct {
	Text string `json:"text"`
}

// TextReview 文字审核
func (l *ReviewClient) TextReview(param TextReviewParam) (ReviewResp, error) {
	resp := ReviewResp{}
	token, _ := l.AccessToken.GetToken()

	u := fmt.Sprintf("%s?access_token=%s", textReviewUrl, token)

	b, _ := json.Marshal(&param)
	body, err := httputil.PostJson(u, b)

	if err != nil {
		logx.Errorf("image review err: %s", err.Error())
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
