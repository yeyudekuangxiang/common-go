package charge

import (
	"encoding/base64"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"strings"
	"testing"
)

var Token = "9e309491-b2f2-4c49-80d6-303074ade1ee"

//开发票1
func TestNotificationMspPaymentInfo(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	status, err := c.NotificationMspPaymentInfo(NotificationMspPaymentInfoParam{
		StartChargeSeq:       "MA1G55M8XufGQgWWL2RSEz5xxlh",
		UserPaidAmount:       0.2,
		SubsidyElecAmount:    0,
		SubsidyServiceAmount: 0,
		RoyaltyServiceAmount: 0.02,
		ChannelAmount:        0,
	})
	println(status)
	if err != nil {
		return
	}
}

//开发票2
func TestInvoiceApply(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	status, err := c.InvoiceApply(InvoiceApplyParam{
		OrderType:    1,
		OutInvoiceId: "MA1G55M8X_test1",
		InvoiceOrders: []InvoiceOrders{
			{
				StartChargeSeq: "MA1G55M8XufGQgWWL2RSEz5xxlh",
				ElecMoney:      0.16,
				SeviceMoney:    0.04,
				TotalMoney:     0.20,
			},
		},
		BusinessType:    2,
		InvoiceTitle:    "颖投信息科技（上海）有限公司",
		InvoiceType:     1,
		PayerRegisterNo: "91310000MA1FY5992T",
		Remark:          "",
		PayerEmail:      "18840853003@163.com",
		ReceiverName:    "刘梅",
		ReceiverAddress: "上海市上海中心",
		ReceiverPhone:   "18840853003",
		CompanyPhone:    "021-62818881",
		CompanyAddress:  "上海市中辉路60号6幢1层137室",
		BankName:        "上海银行静安支行",
		Account:         "03003134635",
	})
	println(status)
	if err != nil {
		return
	}
}

//查看发票

func TestInvoiceInfo(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	status, err := c.InvoiceInfo(InvoiceInfoParam{
		OutInvoiceId: "",
		BatchNo:      "",
		SubBatchNo:   "",
	})
	println(status)
	if err != nil {
		return
	}
}

func TestGetConnectorID(t *testing.T) {
	a := RestoreConnectorID("11000000000000010050638000")
	b := GetConnectorID("https://qrcode.starcharge.com/#/10050638")
	println(b)
	println(a)
	qrCode := "9813627502"
	// 截取终端编号
	if strings.HasPrefix(qrCode, "https://qrcode.starcharge.com/#/") {
		qrCode = strings.TrimPrefix(qrCode, "https://qrcode.starcharge.com/#/")
	}
	// 转换为互联互通编号
	var connectorID string
	if len(qrCode) == 8 {
		connectorID = "110000000000000" + qrCode + "000"
	} else if len(qrCode) == 10 {
		connectorID = "120000000000000" + qrCode[:8] + "0" + qrCode[8:]
	} else {
		return
	}

	fmt.Println("Terminal ID:", qrCode)
	fmt.Println("Connector ID:", connectorID)

}

func TestPhoneEncrypt(t *testing.T) {
	//数据加解
	AESSecret := "OzxlBNxflRPwbePa"
	AESIv := "xnEKN6vfqegWRsbw"
	pkcs5, err := encrypttool.AesEncryptPKCS5([]byte("17601258539"), []byte(AESSecret), []byte(AESIv))
	if err != nil {
		return
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)
	println(data)
}
func TestPhoneDecode(t *testing.T) {
	AESSecret := "OzxlBNxflRPwbePa"
	AESIv := "xnEKN6vfqegWRsbw"
	PhoneEncrypt := "y32i/F+ATATTbqthvk8qyA=="
	decodeString, err := base64.StdEncoding.DecodeString(PhoneEncrypt)
	if err != nil {

	}
	phoneByte, err := encrypttool.AesDecryptPKCS5(decodeString, []byte(AESSecret), []byte(AESIv))

	println(string(phoneByte))
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

//设备验证
func TestQueryEquipAuth(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.QueryEquipAuth(QueryEquipAuthParam{
		EquipAuthSeq: "MA1G55M8X1",
		ConnectorID:  "12000000000000098136272001",
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
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, err := c.QueryStationsInfo(QueryStationsInfoParam{
		LastQueryTime: "",
		PageNo:        1,
		PageSize:      30,
	})
	if err != nil {
		return
	}
	println(resp)

}

//查询充电状态
func TestQueryEquipChargeStatus(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	status, err := c.QueryEquipChargeStatus(QueryEquipChargeStatusParam{
		StartChargeSeq: "MA1G55M8Xh8uMNTmldoYxKazqu2",
	})
	println(status.ConnectorStatus)
	if err != nil {
		return
	}
}

//查询策落
func TestQueryEquipBusinessPolicy(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	status, err := c.QueryEquipBusinessPolicy(QueryEquipBusinessPolicyParam{
		EquipBizSeq: GenerateSerialNumber(),
		ConnectorID: "12000000000000098136275002",
	})
	println(status)
	if err != nil {
		return
	}
}

func TestQueryStationStatus(t *testing.T) {
	c := Client{
		Domain:     "http://test-evcs.starcharge.com/evcs/starcharge/",
		Version:    "",
		AESSecret:  "a2164ada0026ccf7",
		AESIv:      "82c91325e74bef0f",
		SigSecret:  "9af2e7b2d7562ad5",
		Token:      Token,
		OperatorID: "MA1G55M8X",
	}
	status, err := c.QueryStationStatus(QueryStationStatusParam{
		StationIDs: []string{"33221994", "33221993", "33221995", "33221989", "33221990", "89934319"},
	})
	println(status)
	if err != nil {
		return
	}
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
	number := GenerateSerialNumber()
	println(number)
}
