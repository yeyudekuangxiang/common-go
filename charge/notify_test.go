package charge

import (
	"encoding/json"
	"fmt"
	"testing"
)

/**
{OperatorID:313744932 Data:w2F+8JvtXbiaN7aWmWJYuMVwDfGqGAkhdbtz9O+xX/jckUyl3UXdjwA+r6sQaevE5j5+16585wh/Ixs+IAwIBrxYWFgMWvZCQRJlh7/EW5wKW2lM5NOnnr5s5dNmuQ371BRVHgRFAn2fYHtKn2kOoODt2vsJC5NRZjs7L/ekgSRyykFHmuva7W8lTgtvQzxvPjv7UtFwsCXrZSJSoR1mjxUYnUi2ad86e9InekgcaE6c/pbMh2RB9sWxFBEQCl770PrYyXxI7jDv3IE+pplbSXTvMkn7ZNWBuxVksWXam2cyBbvL4bsg2gmMdQleWey3dOQ2q4jL5Fh27lvAPRpszGwEBTmlZDPLOMRNmLQQc3QaoUrgOy8zWvvgZdZjz0TXNmmNIAdtoqwGaGxwLNiqUUw4gb3ei8I4+qOTbI1KRY5kzK/Q9EguaX0bb/hQxvUj7VJZti+8eT/TnYhtZexy70tbKW3qkBIvo3h8lFfw7+kgHJAeFkJRGbrrByy74mkeNuoP3UtBY657KkIPhsQorywAGREGldtOByYBMa/2zPxBXwLzXQGjS63deB86HerGGL5II4veF4IcJ6x4qviK2/fwrp6GVentIGm6vkL41SC/GPBnhv9TNF7lIC0xVvw30DUSOxfhR8vnLPblFWvLfQvIrXyDpt29PSQnGY2EhzM= TimeStamp:20230902184541 Seq:0001 Sig:7685298B0FD0D589E98ED7FCDA6FAA1F}
*/

func TestName(t *testing.T) {
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

//获取绿喵token
func TestNotifyToken(t *testing.T) {
	c := Client{
		Domain:     "https://go-api.miotech.com/mp2c-micro/evcs/v1/",
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

//获取绿喵token
func TestNotifyTokenV2(t *testing.T) {
	c := Client{
		Domain:     "https://go-api.miotech.com/mp2c-micro/evcs/v1/",
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

func TestNotificationInvoiceChange(t *testing.T) {
	c := NotifyClient{
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}
	result := c.QueryRequestEncrypt(NotificationInvoiceChangeParam{
		OutInvoiceId:    "1",
		BatchNo:         "2",
		SubBatchNo:      "3",
		TotalCount:      4,
		InvoiceMaterial: 5,
		InvoiceTime:     "6",
		Status:          7,
		EInvoiceUrl:     "8",
		EInvoiceMiniUrl: "9",
		StartChargeSeqs: "10",
		InvoiceAmount:   11,
		PickupAddress:   "12",
	})
	resultMarshal, err := json.Marshal(result)

	println(resultMarshal)
	resp, err := c.NotificationInvoiceChange(NotificationParam{
		Sig:        "7685298B0FD0D589E98ED7FCDA6FAA1F",
		Data:       "w2F+8JvtXbiaN7aWmWJYuMVwDfGqGAkhdbtz9O+xX/jckUyl3UXdjwA+r6sQaevE5j5+16585wh/Ixs+IAwIBrxYWFgMWvZCQRJlh7/EW5wKW2lM5NOnnr5s5dNmuQ371BRVHgRFAn2fYHtKn2kOoODt2vsJC5NRZjs7L/ekgSRyykFHmuva7W8lTgtvQzxvPjv7UtFwsCXrZSJSoR1mjxUYnUi2ad86e9InekgcaE6c/pbMh2RB9sWxFBEQCl770PrYyXxI7jDv3IE+pplbSXTvMkn7ZNWBuxVksWXam2cyBbvL4bsg2gmMdQleWey3dOQ2q4jL5Fh27lvAPRpszGwEBTmlZDPLOMRNmLQQc3QaoUrgOy8zWvvgZdZjz0TXNmmNIAdtoqwGaGxwLNiqUUw4gb3ei8I4+qOTbI1KRY5kzK/Q9EguaX0bb/hQxvUj7VJZti+8eT/TnYhtZexy70tbKW3qkBIvo3h8lFfw7+kgHJAeFkJRGbrrByy74mkeNuoP3UtBY657KkIPhsQorywAGREGldtOByYBMa/2zPxBXwLzXQGjS63deB86HerGGL5II4veF4IcJ6x4qviK2/fwrp6GVentIGm6vkL41SC/GPBnhv9TNF7lIC0xVvw30DUSOxfhR8vnLPblFWvLfQvIrXyDpt29PSQnGY2EhzM=",
		OperatorID: "313744932",
		TimeStamp:  "20230902184541",
		Seq:        "0001",
	})
	if err != nil {
		return
	}
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(marshal)
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
			ConnectorID: "12000000000000098136271001",
			Status:      2,
			ParkStatus:  50,
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

func TestNotificationStationInfoRequest(t *testing.T) {
	c := NotifyClient{
		AESSecret: "agRigdo8zFu4NMbC",
		AESIv:     "aYqsMbzLCbKpnLLa",
		SigSecret: "dgNaWHDgto716GRd",
	}
	result := c.QueryRequestEncrypt(NotificationStationInfoParam{
		StationInfo: StationInfo{
			StationID:                "33221989",
			OperatorID:               "",
			EquipmentOwnerID:         "1",
			StationName:              "2",
			CountryCode:              "3",
			AreaCode:                 "4",
			Address:                  "5",
			StationTel:               "6",
			ServiceTel:               "7",
			StationType:              8,
			StationStatus:            9,
			ParkNums:                 10,
			StationLng:               11,
			StationLat:               12,
			SiteGuide:                "13",
			Construction:             14,
			Pictures:                 []string{"1212"},
			MatchCars:                "15",
			ParkInfo:                 "16",
			BusineHours:              "17",
			ElectricityFee:           "18",
			ServiceFee:               "19",
			ParkFee:                  "20",
			Payment:                  "21",
			SupportOrder:             22,
			Remark:                   "23",
			ParkingDiscountType:      24,
			ParkFeeStatus:            26,
			BusinessStationFeeDetail: nil,
			StationFeeDetail:         nil,
			OriginalStationFeeDetail: nil,
			IsEnable:                 0,
			PrinterFlag:              0,
			BarrierFlag:              0,
			FloorLevel:               "",
			GuideMap:                 "",
			RoadInfo:                 "",
			AdminName:                "",
			AdminTel:                 "",
			OperationWay:             "",
			EnableRoaming:            0,
			GreenEnergyFlag:          0,
			Flags:                    "",
			EquipmentOperatorID:      "",
			OnlineTime:               "",
			StationGrade:             0,
			SupportingFacilityInfo:   SupportingFacilityInfo{},
			PlaceHolderType:          0,
			EquipmentInfos:           nil,
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
