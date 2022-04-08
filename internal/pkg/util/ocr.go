package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type OCRResult struct {
	WordsResult []struct {
		Words string `json:"words"`
	} `json:"words_result"`
	WordsResultNum int   `json:"words_result_num"`
	LogID          int64 `json:"log_id"`
}

func OCRPush(src string) *OCRResult {

	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage?access_token=24.6157c4c9729181acc1bac04d6bd5ecbe.2592000.1650680140.282335-25833266"
	method := "POST"

	payload := strings.NewReader("url=" + src)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var o OCRResult
	err = json.Unmarshal([]byte(string(body)), &o)
	fmt.Println("OCRPush ", payload, string(body))

	return &o

}
