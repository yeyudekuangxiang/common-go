package api

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/util/encrypt"
	"mio/pkg/wxapp"
	"testing"
	"time"
)

func TestCheckTime(t *testing.T) {
	startTime, _ := time.Parse("2006-01-02", "2022-08-22")
	endTime, _ := time.Parse("2006-01-02", "2022-08-31")
	//fmt.Printf("start error : %s", err.Error())
	fmt.Println(time.Now())
	fmt.Println(startTime)
	fmt.Println(endTime)

	if time.Now().After(startTime) && time.Now().Before(endTime) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func TestMakeSign(t *testing.T) {

	data := RecycleFmyForm{
		AppId:          "75133417",
		NotificationAt: "2022-09-06 11:06:32",
		Sign:           "02871058a4bfde2bffe88d6366865",
		Data: RecycleFmyData{
			OrderSn: "20220906110237856327",
			Status:  "COMPLETE",
			Weight:  "20.00",
			Reason:  "SF1347969593527",
			Phone:   "18301939833",
		},
	}

	rand1 := "0287"
	rand2 := "6865"
	appId := "75133417"
	appSecret := "jlLV7gw9s7xTShQIriSUiyQu5FfBfSZv"

	marshal, err := json.Marshal(data)
	if err != nil {
		return
	}
	jsonData := string(marshal)

	fmt.Println(jsonData)

	verifyData := rand1 + appId + jsonData + appSecret + rand2
	md5Data := encrypt.Md5(verifyData)
	sign := rand1 + string([]rune(md5Data)[7:21]) + rand2
	fmt.Println(sign)

}

func TestAccessToken(t *testing.T) {
	initialize.Initialize("/Users/yunfeng/Documents/workspace/mp2c-go/config.ini")
	token, _ := wxapp.NewClient(app.Weapp).AccessToken()
	fmt.Println(token)
}
