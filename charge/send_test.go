package charge

import (
	"encoding/base64"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"testing"
)

func TestGenerateSerialNumber(t *testing.T) {
	number := GenerateSerialNumber("MA1G55M8X")
	println(number)
}
func TestNotificationStationStatusRequest(t *testing.T) {
	c := Client{
		Domain:       "127.0.0.1:1017/evcs/v1",
		Version:      "",
		AESSecret:    "agRigdo8zFu4NMbC",
		AESIv:        "aYqsMbzLCbKpnLLa",
		SigSecret:    "dgNaWHDgto716GRd",
		Token:        "",
		OperatorID:   "313744932",
		MIOAESSecret: "agRigdo8zFu4NMbC",
		MIOAESIv:     "aYqsMbzLCbKpnLLa",
		MIOSigSecret: "dgNaWHDgto716GRd",
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
	c := Client{
		Domain:       "127.0.0.1:1017/evcs/v1",
		Version:      "",
		AESSecret:    "agRigdo8zFu4NMbC",
		AESIv:        "aYqsMbzLCbKpnLLa",
		SigSecret:    "dgNaWHDgto716GRd",
		Token:        "",
		OperatorID:   "313744932",
		MIOAESSecret: "agRigdo8zFu4NMbC",
		MIOAESIv:     "aYqsMbzLCbKpnLLa",
		MIOSigSecret: "dgNaWHDgto716GRd",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()

	a := NotificationStartChargeResult{
		StartChargeSeq: "1212",
		SuccStat:       1,
		FailReason:     2,
	}
	resp := c.NotificationResponse(a)
	println(resp)
}
func TestQueryEquipAuth(t *testing.T) {
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "82c91325e74bef0f",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "35a04fd7-0d76-43ea-b621-157c4dd2dc12",
		OperatorID:   "MA1G55M8X",
		MIOAESSecret: "",
		MIOAESIv:     "",
		MIOSigSecret: "",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.QueryEquipAuth(QueryEquipAuthParam{
		EquipBizSeq: "MA1G55M8X1",
		ConnectorID: "12000000000000072155462002",
	})
	if err != nil {
		return
	}
	println(resp)
}
func TestQueryStationsInfo(t *testing.T) {
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "82c91325e74bef0f",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "35a04fd7-0d76-43ea-b621-157c4dd2dc12",
		OperatorID:   "MA1G55M8X",
		MIOAESSecret: "",
		MIOAESIv:     "",
		MIOSigSecret: "",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.QueryStationsInfo(QueryStationsInfoParam{
		LastQueryTime: "2023-06-11 01:25:11",
		PageNo:        1,
		PageSize:      1,
	})
	if err != nil {
		return
	}
	println(resp)

}

func TestQueryTokenRequest(t *testing.T) {
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "82c91325e74bef0f",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "35a04fd7-0d76-43ea-b621-157c4dd2dc12",
		OperatorID:   "MA1G55M8X",
		MIOAESSecret: "",
		MIOAESIv:     "",
		MIOSigSecret: "",
	}
	request, err := c.QueryTokenRequest(NotificationParam{
		Sig:        "",
		Data:       "",
		OperatorID: "",
		TimeStamp:  "",
		Seq:        "",
	})
	println(request)
	if err != nil {
		return
	}
}

//星星充电回调用绿喵测试

func TestTokenV2(t *testing.T) {
	c := Client{
		Domain:       "127.0.0.1:1017/evcs/v1",
		Version:      "",
		AESSecret:    "agRigdo8zFu4NMbC",
		AESIv:        "aYqsMbzLCbKpnLLa",
		SigSecret:    "dgNaWHDgto716GRd",
		Token:        "",
		OperatorID:   "313744932",
		MIOAESSecret: "agRigdo8zFu4NMbC",
		MIOAESIv:     "aYqsMbzLCbKpnLLa",
		MIOSigSecret: "dgNaWHDgto716GRd",
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

func TestToken(t *testing.T) {
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "82c91325e74bef0f",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "",
		OperatorID:   "MA1G55M8X",
		MIOAESSecret: "",
		MIOAESIv:     "",
		MIOSigSecret: "",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.QueryToken(QueryTokenParam{
		OperatorSecret: "acb93539fc9bg78k",
		OperatorID:     "MA1G55M8X",
	})
	if err != nil {
		return
	}
	println(resp)

}

func TestNotificationRequest(t *testing.T) {
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "",
		OperatorID:   "",
		MIOAESSecret: "agRigdo8zFu4NMbC",
		MIOAESIv:     "aYqsMbzLCbKpnLLa",
		MIOSigSecret: "dgNaWHDgto716GRd",
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
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "",
		OperatorID:   "",
		MIOAESSecret: "agRigdo8zFu4NMbC",
		MIOAESIv:     "aYqsMbzLCbKpnLLa",
		MIOSigSecret: "dgNaWHDgto716GRd",
	}

	res := c.NotificationResult(QueryResponse{
		Ret:  500,
		Msg:  "错误信息",
		Data: []byte("123"),
		Sig:  "",
	})
	println(res)
}
func TestBikeCard(t *testing.T) {

	//数据加解
	pkcs5, err := encrypttool.AesEncryptPKCS5([]byte("123456"), []byte("agRigdo8zFu4NMbC"), []byte("aYqsMbzLCbKpnLLa"))
	if err != nil {
		return
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)

	println(data)

	decryptPKCS5, err := encrypttool.AesDecryptPKCS5([]byte(pkcs5), []byte("agRigdo8zFu4NMbC"), []byte("aYqsMbzLCbKpnLLa"))
	if err != nil {
		return
	}

	println(decryptPKCS5)

	/*
		println(data)
		c := Client{
			Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
			Version:    "",
			AESSecret:  "a2164ada0026ccf7",
			SigSecret:  "9af2e7b2d7562ad5",
			Token:      "",
			OperatorID: "",
		}
		//bizId := time.Now().Format("20060102150405") + c.rand()
		resp, _ := c.QueryEquipAuth(QueryEquipAuthParam{
			EquipBizSeq: "",
			ConnectorID: "",
		})
	*/
	/*	println(resp.EquipAuthSeq)*/

}
