package platform

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/duiba/util"
	"net/url"
	"testing"
)

func TestGetToken(t *testing.T) {
	signStr := sign()
	//params := struct {
	//	AppId    string `json:"appId"`
	//	ClientId string `json:"clientId"`
	//	Sign     string `json:"sign"`
	//}{
	//	AppId:    "8c216e9170a1426d95621ba511d34bb5",
	//	ClientId: "o0GhT6mOSzL9Nm0g2VrQD8QPTps8",
	//	Sign:     signStr,
	//}
	params := make(url.Values)
	params.Set("appId", "8c216e9170a1426d95621ba511d34bb5")
	params.Set("clientId", "o0GhT6mOSzL9Nm0g2VrQD8QPTps8")
	params.Set("sign", signStr)
	u := "https://uat.oola.cn/oola-api/api/user/getUserAutoLoginKey"
	body, err := httputil.PostFrom(u, params)
	fmt.Printf("body:%s\n", body)
	if err != nil {
		fmt.Printf("http error:%s\n", err.Error())
	}
	//response
	res := oolaResponse{}
	if err = json.Unmarshal(body, &res); err != nil {
		fmt.Printf("json error:%s\n", err.Error())
	}
	fmt.Printf("response: %v\n", res)
}

func sign() string {
	return util.Md5("83712eb73e794e569b8565f61053dc05appId=8c216e9170a1426d95621ba511d34bb5;clientId=o0GhT6mOSzL9Nm0g2VrQD8QPTps8;")
}

func TestRegister(t *testing.T) {
	signStr := sign()
	params := make(url.Values)
	params.Set("appId", "8c216e9170a1426d95621ba511d34bb5")
	params.Set("clientId", "o0GhT6mOSzL9Nm0g2VrQD8QPTps8")
	params.Set("sign", signStr)
	u := "https://uat.oola.cn/oola-api/api/user/register"
	body, err := httputil.PostFrom(u, params)
	fmt.Printf("body:%s\n", body)
	if err != nil {
		fmt.Printf("http error:%s\n", err.Error())
	}
	//response
	res := oolaResponse{}
	if err = json.Unmarshal(body, &res); err != nil {
		fmt.Printf("json error:%s\n", err.Error())
	}
	fmt.Printf("response: %v\n", res)
}

func TestLoginKey(t *testing.T) {
	signStr := sign()
	params := make(url.Values)
	params.Set("appId", "8c216e9170a1426d95621ba511d34bb5")
	params.Set("clientId", "o0GhT6mOSzL9Nm0g2VrQD8QPTps8")
	params.Set("sign", signStr)
	u := "https://uat.oola.cn/oola-api/api/user/getUserAutoLoginKey"
	body, err := httputil.PostFrom(u, params)
	fmt.Printf("body:%s\n", body)
	if err != nil {
		fmt.Printf("http error:%s\n", err.Error())
	}
	res := oolaResponse{}
	if err = json.Unmarshal(body, &res); err != nil {
		fmt.Printf("json error:%s\n", err.Error())
	}
	fmt.Printf("response: %v\n", res)
}
