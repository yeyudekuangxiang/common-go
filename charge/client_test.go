package charge

import (
	"encoding/base64"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"testing"
)

//设备验证
func TestQueryEquipAuth(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      "35a04fd7-0d76-43ea-b621-157c4dd2dc12",
		OperatorID: "MA1G55M8X",
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

//拉去充电站信息
func TestQueryStationsInfo(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      "35a04fd7-0d76-43ea-b621-157c4dd2dc12",
		OperatorID: "MA1G55M8X",
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

//获取星星token
func TestToken(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      "",
		OperatorID: "MA1G55M8X",
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

func TestGenerateSerialNumber(t *testing.T) {
	number := GenerateSerialNumber("MA1G55M8X")
	println(number)
}
