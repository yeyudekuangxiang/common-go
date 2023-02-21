package unidian

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"io/ioutil"
	"net/http"
	"time"
)

func CouponOfUnidian(typeId string, mobile string, outTradeNo string) {

	timeStamp := time.Now().Format("20060102150405")

	channelId := "115"
	//timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	key := "B5C0B2C3C1CD4942"
	sign := encrypttool.Md5(typeId + "#" + channelId + "#" + mobile + "#" + timeStamp + "#" + outTradeNo + "#" + key)
	//sign := util.Md5(typeId + "#" + channelId + "#" + mobile + "#" + timeStamp + "#" + outTradeNo + "#" + key)
	url := "http://qyif.unidian.com/QuanYi/Common/Coupon.aspx?TypeId=" + typeId + "&ChannelId=" + channelId + "&Mobile=" + mobile + "&TimeStamp=" + timeStamp + "&OutTradeNo=" + outTradeNo + "&Sign=" + sign + "&UserIdType=0"
	method := "GET"

	fmt.Println(url)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "ASP.NET_SessionId=x3ccfx031odfew5d4irdsy0b")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
