package charge

import (
	"encoding/json"
	"fmt"
	"testing"
)

//获取绿喵token
func TestNotifyToken(t *testing.T) {
	c := Client{
		Domain:     "https://godev-api.miotech.com/mp2c-micro/evcs/v1/",
		Version:    "",
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		Token:      "",
		OperatorID: "313744932",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.QueryToken(QueryTokenParam{
		OperatorSecret: "NU0gYnwsQaLTAQ0loRwol4NaRx8tZksX",
		OperatorID:     "313744932",
	})
	if err != nil {
		return
	}
	println(resp)
}

func TestNotificationStartChargeResultParams(t *testing.T) {
	c := NotifyClient{
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		OperatorID: "313744932",
	}
	result := c.QueryRequestEncrypt(NotificationStartChargeResultParam{
		StartChargeSeq:     "MA1G55M8X633322921",
		StartChargeSeqStat: 1,
		ConnectorID:        "12000000000000072155475002",
		StartTime:          "2023-06-13 18:02:37",
		IdentCode:          "IdentCode",
	})
	resultMarshal, err := json.Marshal(result)

	println(resultMarshal)
	resp, err := c.NotificationStartChargeResultRequest(*result)
	if err != nil {
		return
	}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(marshal)
}

func TestNotificationEquipChargeStatus(t *testing.T) {
	c := NotifyClient{
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		OperatorID: "313744932",
	}
	result := c.QueryRequestEncrypt(NotificationEquipChargeStatusParam{
		StartChargeSeq:     "MA1G55M8X633322921",
		StartChargeSeqStat: 1,
		ConnectorID:        "1212",
		StartTime:          "2023-06-13 20:40:05",
		IdentCode:          "",
		ConnectorStatus:    1,
		CurrentA:           0,
		CurrentB:           0,
		CurrentC:           0,
		VoltageA:           0,
		VoltageB:           0,
		VoltageC:           0,
		Soc:                0,
		EndTime:            "2023-06-13 20:40:05",
		TotalPower:         0,
		ElecMoney:          0,
		ServiceMoney:       0,
		TotalMoney:         0,
		StopReason:         0,
		SumPeriod:          0,
		ChargeDetails:      "",
	})
	resultMarshal, err := json.Marshal(result)

	println(resultMarshal)
	resp, err := c.NotificationStartChargeResultRequest(*result)
	if err != nil {
		return
	}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(marshal)
}

func TestNotificationStationStatusRequest(t *testing.T) {
	c := NotifyClient{
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		OperatorID: "313744932",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.NotificationStationStatusRequest(NotificationParam{
		Sig:        "",
		Data:       "",
		OperatorID: "",
		TimeStamp:  "",
		Seq:        "",
	})
	if err != nil {
		return
	}
	println(resp)
}

type NotificationStartChargeResult struct {
	StartChargeSeq string `json:"StartChargeSeq"`
	SuccStat       int64  `json:"SuccStat"`
	FailReason     int64  `json:"failReason"`
}

func TestNotificationResponse(t *testing.T) {
	c := NotifyClient{
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		OperatorID: "313744932",
	}
	a := NotificationStartChargeResult{
		StartChargeSeq: "1212",
		SuccStat:       1,
		FailReason:     2,
	}
	resp := c.NotificationResponse(a)
	println(resp)
}

func TestNotificationRequest(t *testing.T) {
	c := NotifyClient{
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		OperatorID: "313744932",
	}
	validate, err := c.NotificationRequest(NotificationParam{
		Sig:        "E0E972AB13A63F38D6B228FE656FB5DE",
		Data:       "NJJ5Fk6xAcU8d6lpqQhmPg==",
		OperatorID: "1212",
		TimeStamp:  "1212",
		Seq:        "121212",
	})
	if err != nil {
		return
	}
	println(validate)
}

func TestNotificationResult(t *testing.T) {
	c := NotifyClient{
		AESSecret:  "agRigdo8zFu4NMbC",
		AESIv:      "aYqsMbzLCbKpnLLa",
		SigSecret:  "dgNaWHDgto716GRd",
		OperatorID: "313744932",
	}

	res := c.NotificationResult(QueryResponse{
		Ret:  500,
		Msg:  "错误信息",
		Data: []byte("123"),
		Sig:  "",
	})
	println(res)
}
