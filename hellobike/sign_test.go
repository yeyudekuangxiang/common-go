package hellobike

import (
	"strconv"
	"testing"
	"time"
)

//{"orderNo":"TP20230317181438200002103590028"}

func TestRefundCard(t *testing.T) {
	c := Client{
		AppId:   "20230302145050102",
		Version: "1.0",
		Action:  "hellobike.tw.refundcard",
		AppKey:  "d9244321dc3246caa54a29e7c156dd0c",
		Domain:  "https://openapi.hellobike.com/open/api",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	//202303171645544129841576
	resp, _ := c.RefundCard(RefundCardParam{
		AppId:     c.AppId,
		Action:    c.Action,
		Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Sign:      "",
		Version:   c.Version,
		BizContent: struct {
			ActivityId    string `json:"activityId"`
			OrderNo       string `json:"orderNo"`
			MobilePhone   string `json:"mobilePhone"`
			TransactionId string `json:"transactionId"`
		}{
			ActivityId:    "H3979885952972083867",
			OrderNo:       "TP20230317181438200002103590028",
			MobilePhone:   "18840853003",
			TransactionId: "",
		},
	})
	println(resp.Code)

}

func TestBikeCard(t *testing.T) {
	c := Client{
		AppId:   "20230302145050102",
		Version: "1.0",
		Action:  "hellobike.activity.bikecard",
		AppKey:  "d9244321dc3246caa54a29e7c156dd0c",
		Domain:  "https://openapi.hellobike.com/bike/activity",
	}

	bizId := time.Now().Format("20060102150405") + c.rand()
	//202303171645544129841576
	resp, _ := c.SendCoupon(SendCouponParam{
		AppId:        c.AppId,
		Action:       c.Action,
		UtcTimestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Sign:         "",
		Version:      c.Version,
		BizContent: struct {
			ActivityId    string `json:"activityId"`
			MobilePhone   string `json:"mobilePhone"`
			TransactionId string `json:"transactionId"`
		}{
			ActivityId:    "H3979885952972083867",
			MobilePhone:   "18840853003",
			TransactionId: bizId,
		},
	})
	println(resp.ErrorCode)

}
func TestHelloBike(t *testing.T) {
	c := Client{
		AppId:   "20200907153742407",
		Version: "1.0",
		Action:  "hellobike.activity.bikecard",
		AppKey:  "75e3747b359246379b2447dfd5090b8a",
	}
	resp, _ := c.SendCoupon(SendCouponParam{
		AppId:        c.AppId,
		Action:       c.Action,
		UtcTimestamp: "1599634041750",
		Sign:         "",
		Version:      c.Version,
		BizContent: struct {
			ActivityId    string `json:"activityId"`
			MobilePhone   string `json:"mobilePhone"`
			TransactionId string `json:"transactionId"`
		}{
			ActivityId:    "1296103963453030400",
			MobilePhone:   "13661502232",
			TransactionId: "202009091447R0001243",
		},
	})
	println(resp.ErrorCode)
}
