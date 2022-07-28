package baidu

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/util/httputil"
)

const (
	webimageUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage"
)

type ImageClient struct {
	AccessToken *AccessToken
}
type WebImageParam struct {
	ImageUrl string
}
type WebImageResult struct {
	ErrorResponse
	LogId          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

func (c ImageClient) WebImage(param WebImageParam) (*WebImageResult, error) {
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
