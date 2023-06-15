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
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
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
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}
	ChargeDetailsList := make([]ChargeDetails, 0)
	ChargeDetailsList = append(ChargeDetailsList, ChargeDetails{
		DetailPower:       11.4,
		ElecPrice:         0.8,
		SevicePrice:       0.2,
		DetailElecMoney:   9.12,
		DetailSeviceMoney: 2.28,
		DetailStartTime:   "2023-06-14 16:43:00",
		DetailEndTime:     "2023-06-14 17:59:00",
		DetailType:        2,
	})
	result := c.QueryRequestEncrypt(NotificationEquipChargeStatusParam{
		StartChargeSeq:     "MA1G55M8XtmRTO6reuZLoiwJwq0",
		StartChargeSeqStat: 2,
		ConnectorID:        "12000000000000098136275002",
		ConnectorStatus:    3,
		CurrentA:           454.55,
		CurrentB:           0,
		CurrentC:           0,
		VoltageA:           220,
		VoltageB:           0,
		VoltageC:           0,
		Soc:                0,
		StartTime:          "2023-06-14 16:43:50",
		EndTime:            "2023-06-14 17:59:39",
		TotalPower:         11.4,
		ElecMoney:          9.12,
		SeviceMoney:        2.28,
		TotalMoney:         11.4,
		SumPeriod:          1,
		ChargeDetails:      ChargeDetailsList,
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

func TestNotificationStopChargeResultRequest(t *testing.T) {
	c := NotifyClient{
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}
	result := c.QueryRequestEncrypt(NotificationStopChargeResultParam{
		StartChargeSeq:     "313744932651370500",
		StartChargeSeqStat: 4,
		ConnectorID:        "12000000000000072155475002",
		SuccStat:           0,
		FailReason:         0,
	})
	resultMarshal, err := json.Marshal(result)

	println(resultMarshal)
	resp, err := c.NotificationStopChargeResultRequest(*result)
	if err != nil {
		return
	}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(marshal)
}

func TestNNotificationChargeOrderInfoRequest(t *testing.T) {
	c := NotifyClient{
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}
	ChargeDetail := make([]ChargeDetails, 0)

	ChargeDetail = append(ChargeDetail, ChargeDetails{
		DetailPower:       0,
		ElecPrice:         0,
		SevicePrice:       0,
		DetailElecMoney:   0,
		DetailSeviceMoney: 0,
		DetailStartTime:   "1212",
		DetailEndTime:     "121212",
		DetailType:        0,
	})
	result := c.QueryRequestEncrypt(NotificationChargeOrderInfoParam{
		StartChargeSeq:   "MA1G55M8X633322921",
		ConnectorId:      "12",
		StartTime:        "2023-06-14 00:06:20",
		EndTime:          "2023-06-14 00:06:20",
		TotalPower:       10,
		TotalElecMoney:   11,
		TotalSeviceMoney: 12,
		TotalMoney:       13,
		StopReason:       0,
		SumPeriod:        0,
		ChargeDetails:    ChargeDetail,
	})
	resultMarshal, err := json.Marshal(result)

	println(resultMarshal)
	resp, err := c.NotificationChargeOrderInfoRequest(*result)
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
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}
	result := c.QueryRequestEncrypt(NotificationStationStatusParam{
		StationStatusInfo: ConnectorStatusInfo{
			ConnectorID: "12000000000000098136275002",
			Status:      1,
			ParkStatus:  1,
			LockStatus:  1,
			StatusCode:  "12121212",
		},
	})
	resultMarshal, err := json.Marshal(result)

	println(resultMarshal)
	resp, err := c.NotificationStopChargeResultRequest(*result)
	if err != nil {
		return
	}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(marshal)
}

type NotificationStartChargeResult struct {
	StartChargeSeq string `json:"StartChargeSeq"`
	SuccStat       int64  `json:"SuccStat"`
	FailReason     int64  `json:"failReason"`
}

func TestNotificationResponse(t *testing.T) {
	c := NotifyClient{
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
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
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
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
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}

	res := c.NotificationResult(QueryResponse{
		Ret:  500,
		Msg:  "错误信息",
		Data: []byte("123"),
		Sig:  "",
	})
	println(res)
}
