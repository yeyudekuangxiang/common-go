package api

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/util/encrypt"
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
		AppId:          "52558668",
		NotificationAt: "2022-09-05 16:26:57",
		//Sign:           "1524196360985d37a99a0df044613",
		Data: RecycleFmyData{
			OrderSn: "20220905161458561515",
			Status:  "COMPLETE",
			Weight:  "10.00",
			Reason:  "SF1347643098218",
			Phone:   "18301939833",
		},
	}
	rand1 := "1524"
	rand2 := "4613"
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
