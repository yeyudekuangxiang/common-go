package wxamp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RiskCase struct {
	BatchID  string `json:"batchId"`
	RecordID string `json:"recordId"`
	RPCCount int    `json:"__rpcCount"`
}

type CaseRes struct {
	List []struct {
		Openid    string   `json:"openid"`
		RiskRank  int      `json:"riskRank"`
		Scene     int      `json:"scene"`
		Errcode   int      `json:"errcode"`
		LabelList []string `json:"labelList"`
	} `json:"list"`
	RPCCount int `json:"__rpcCount"`
}

func BatchGetUserRiskCase(openidArr []string) *CaseRes {
	url := "https://mp.weixin.qq.com/wxamp/cgi/development/BatchGetUserRiskCase?token=523498074&lang=zh_CN&random=0.0736900020659983"
	method := "POST"
	openids := ""
	for _, i := range openidArr {
		openids += "%22" + i + "%22%2C"
	}
	openids = openids[0 : len(openids)-3]
	payload := strings.NewReader("data=%7B%22openidList%22%3A%5B" + openids + "%5D%2C%22fileOpenidList%22%3A%5B%5D%2C%22sceneList%22%3A%5B0%5D%7D")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("authority", "mp.weixin.qq.com")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", "ua_id=FoTFMNyH1cSbFoa5AAAAANxMnJ91Zbrnb3Gbn80y-UI=; wxuin=52349754594085; mm_lang=zh_CN; utm_source=baidu; utm_medium=sem; _ga=GA1.2.774026700.1652667909; _hjSessionUser_2765497=eyJpZCI6IjM4MmEwZTVmLWIzYjMtNWRkYS1iMWVjLWI3OTM1MDVlZDg4NiIsImNyZWF0ZWQiOjE2NTI2Njc5MTA1MjEsImV4aXN0aW5nIjp0cnVlfQ==; uuid=9b476ab24f332dc61ef7a330b3a37369; rand_info=CAESIJ5isIhQgAxG3xsrgsO/IaKPhHro2g1/rvMF6d+4FnuR; slave_bizuin=3829461195; data_bizuin=3829461195; bizuin=3829461195; data_ticket=NB3BFsfkIQJL7V4xXSCZBC4L7dlzww4Q0A58btAE6XhBKL1EV85QhbwVpZhvuArw; slave_sid=WnlNd3RfWVpVWXhjU0NCY1NJbl93SXV6dENHYTRoYkxjUUtQWElmRFBod0VsWDB4WEhTTWpGREdBeWYyY0JQS1RkNlI5emZ4ZXJoNThEbHR0NGFWZnZhRFFMaUM3OGttWTRoTDNiX1FTcnN4Qmg4dkJSVnNZOGtJUXF0NEtmaDN3VG54b2Y5MjM1VmRYN3VG; slave_user=gh_45da3fcaeadd; xid=85c2b602859675376d5cf1b228389b81; sig=h014b58b8e895d77bdbc053c9a74eae295ce265b2cea218b0f6d9395597edf65c6ca96457a162d0d4a7")
	req.Header.Add("origin", "https://mp.weixin.qq.com")
	req.Header.Add("referer", "https://mp.weixin.qq.com/wxamp/security/user?token=523498074&lang=zh_CN")
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var c RiskCase
	err = json.Unmarshal([]byte(string(body)), &c)
	fmt.Println("RiskCase ", payload, string(body))

	r := GetUserRiskCaseResult(c.BatchID)
	return r
}

func GetUserRiskCaseResult(batchId string) *CaseRes {
	url := "https://mp.weixin.qq.com/wxamp/cgi/development/GetUserRiskCaseResult?token=523498074&lang=zh_CN&random=0.4714267962223181"
	method := "POST"

	payload := strings.NewReader("batchId=" + batchId)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("authority", "mp.weixin.qq.com")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", "ua_id=FoTFMNyH1cSbFoa5AAAAANxMnJ91Zbrnb3Gbn80y-UI=; wxuin=52349754594085; mm_lang=zh_CN; utm_source=baidu; utm_medium=sem; _ga=GA1.2.774026700.1652667909; _hjSessionUser_2765497=eyJpZCI6IjM4MmEwZTVmLWIzYjMtNWRkYS1iMWVjLWI3OTM1MDVlZDg4NiIsImNyZWF0ZWQiOjE2NTI2Njc5MTA1MjEsImV4aXN0aW5nIjp0cnVlfQ==; uuid=9b476ab24f332dc61ef7a330b3a37369; rand_info=CAESIJ5isIhQgAxG3xsrgsO/IaKPhHro2g1/rvMF6d+4FnuR; slave_bizuin=3829461195; data_bizuin=3829461195; bizuin=3829461195; data_ticket=NB3BFsfkIQJL7V4xXSCZBC4L7dlzww4Q0A58btAE6XhBKL1EV85QhbwVpZhvuArw; slave_sid=WnlNd3RfWVpVWXhjU0NCY1NJbl93SXV6dENHYTRoYkxjUUtQWElmRFBod0VsWDB4WEhTTWpGREdBeWYyY0JQS1RkNlI5emZ4ZXJoNThEbHR0NGFWZnZhRFFMaUM3OGttWTRoTDNiX1FTcnN4Qmg4dkJSVnNZOGtJUXF0NEtmaDN3VG54b2Y5MjM1VmRYN3VG; slave_user=gh_45da3fcaeadd; xid=85c2b602859675376d5cf1b228389b81; sig=h014b58b8e895d77bdbc053c9a74eae295ce265b2cea218b0f6d9395597edf65c6ca96457a162d0d4a7")
	req.Header.Add("origin", "https://mp.weixin.qq.com")
	req.Header.Add("referer", "https://mp.weixin.qq.com/wxamp/security/user?token=523498074&lang=zh_CN")
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var r CaseRes
	err = json.Unmarshal([]byte(string(body)), &r)
	fmt.Println("RiskCase ", payload, string(body))

	return &r
}
