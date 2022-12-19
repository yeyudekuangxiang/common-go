package baidu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mio/config"
	"net/http"
	"net/url"
	"strings"
)

type ReviewClient struct {
	ApiKey    string
	SecretKey string
	image     string
}

func NewReviewClient() *ReviewClient {
	return &ReviewClient{
		ApiKey:    config.Config.BaiDuReview.AppKey,
		SecretKey: config.Config.BaiDuReview.AppSecret,
	}
}

// LoadImageUrl 参数为图片的url地址，要求图片base64后大于5kb小于4m
func (l *ReviewClient) LoadImageUrl(imageUrl string) {
	l.image = url.QueryEscape(imageUrl)
}

// LoadImage 参数为base64后的图片，要求图片base64后大于5kb小于4m
func (l *ReviewClient) LoadImage(image string) {
	l.image = image
}

// Review 审核
func (l *ReviewClient) Review() ReviewResp {
	resp := ReviewResp{}
	if l.image == "" {
		return resp
	}
	u := "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_defined?access_token=" + l.GetAccessToken()
	payload := strings.NewReader(l.image)
	client := &http.Client{}
	req, err := http.NewRequest("POST", u, payload)

	if err != nil {
		fmt.Println(err)
		return resp
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return resp
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return resp
	}
	fmt.Println(string(body))
	json.Unmarshal(body, &resp)
	fmt.Println(resp)
	return resp
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func (l *ReviewClient) GetAccessToken() string {
	u := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", l.ApiKey, l.SecretKey)
	resp, err := http.Post(u, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	accessTokenObj := map[string]string{}
	json.Unmarshal([]byte(body), &accessTokenObj)
	return accessTokenObj["access_token"]
}

//response
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
