package charge

import (
	"encoding/base64"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"testing"
)

//星星充电回调用绿喵测试

func TestToken(t *testing.T) {
	c := Client{
		Domain:       "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:      "",
		AESSecret:    "a2164ada0026ccf7",
		AESIv:        "82c91325e74bef0f",
		SigSecret:    "9af2e7b2d7562ad5",
		Token:        "",
		OperatorID:   "MA1G55M8X1",
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

	res := c.NotificationResult(starChargeResponse{
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
