package charge

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"strconv"
	"strings"
	"time"
)

type NotifyClient struct {
	AESSecret  string //acb93539fc9bg78k
	AESIv      string //82c91325e74bef0f
	SigSecret  string //9af2e7b2d7562ad5  签名密钥
	Token      string
	OperatorID string
}

//请求绿喵token

func (c *NotifyClient) QueryTokenRequest(param NotificationParam) (resp *QueryTokenReq, err error) {
	request, err := c.NotificationRequest(param)
	if err != nil {
		return nil, err
	}
	QueryToken := QueryTokenReq{}
	err = json.Unmarshal(request, &QueryToken)
	if err != nil {
		return nil, err
	}
	return &QueryToken, nil
}

//星星回调绿喵接口

func (c *NotifyClient) NotificationRequest(param NotificationParam) (resp []byte, err error) {
	operatorID := param.OperatorID
	data := param.Data
	timestamp := param.TimeStamp
	seq := param.Seq

	encReq := operatorID + data + timestamp + seq
	sign := strings.ToUpper(encrypttool.HMacMd5(encReq, c.SigSecret))

	if sign != param.Sig {
		return nil, errors.New("签名失败")
	}
	dataString, err := base64.StdEncoding.DecodeString(param.Data)
	if err != nil {
		return nil, err
	}
	//数据解密
	decryptPKCS5, err := encrypttool.AesDecryptPKCS5(dataString, []byte(c.AESSecret), []byte(c.AESIv))
	if err != nil {
		return nil, err
	}
	return decryptPKCS5, nil
}

//返回结果加密返回

func (c *NotifyClient) NotificationResult(param QueryResponse) (resp *ChargeResponse) {
	if param.Ret != 0 {
		return &ChargeResponse{
			Ret:  param.Ret,
			Msg:  param.Msg,
			Data: "",
			Sig:  "",
		}
	}
	pkcs5, err := encrypttool.AesEncryptPKCS5(param.Data, []byte(c.AESSecret), []byte(c.AESIv))
	if err != nil {
		return &ChargeResponse{
			Ret:  param.Ret,
			Msg:  err.Error(),
			Data: "",
			Sig:  "",
		}
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)
	//返回值加签
	encResp := strconv.Itoa(param.Ret) + InterfaceToString(param.Msg) + data
	sign := strings.ToUpper(encrypttool.HMacMd5(encResp, c.SigSecret))
	res := ChargeResponse{
		Ret:  param.Ret,
		Msg:  param.Msg,
		Data: data,
		Sig:  sign,
	}
	return &res
}

//设备状态变化推送

func (c *NotifyClient) NotificationStationStatusRequest(param NotificationParam) (resp *NotificationStationStatusParam, err error) {
	request, err := c.NotificationRequest(param)
	if err != nil {
		return nil, err
	}
	stationStatus := NotificationStationStatusParam{}
	err = json.Unmarshal(request, &stationStatus)
	if err != nil {
		return nil, err
	}
	return &stationStatus, nil
}

//推送充电结果

func (c *NotifyClient) NotificationStartChargeResultRequest(param NotificationParam) (resp *NotificationStartChargeResultParam, err error) {
	request, err := c.NotificationRequest(param)
	if err != nil {
		return nil, err
	}
	stationStatus := NotificationStartChargeResultParam{}
	err = json.Unmarshal(request, &stationStatus)
	if err != nil {
		return nil, err
	}
	return &stationStatus, nil
}

//推送充电状态

func (c *NotifyClient) NotificationEquipChargeStatusRequest(param NotificationParam) (resp *NotificationEquipChargeStatusParam, err error) {
	request, err := c.NotificationRequest(param)
	if err != nil {
		return nil, err
	}
	stationStatus := NotificationEquipChargeStatusParam{}
	err = json.Unmarshal(request, &stationStatus)
	if err != nil {
		return nil, err
	}
	return &stationStatus, nil
}

//推送停止充电结果

func (c *NotifyClient) NotificationStopChargeResultRequest(param NotificationParam) (resp *NotificationStopChargeResultParam, err error) {
	request, err := c.NotificationRequest(param)
	if err != nil {
		return nil, err
	}
	stationStatus := NotificationStopChargeResultParam{}
	err = json.Unmarshal(request, &stationStatus)
	if err != nil {
		return nil, err
	}
	return &stationStatus, nil
}

//推送充电订单信息

func (c *NotifyClient) NotificationChargeOrderInfoRequest(param NotificationParam) (resp *NotificationChargeOrderInfoParam, err error) {
	request, err := c.NotificationRequest(param)
	if err != nil {
		return nil, err
	}
	stationStatus := NotificationChargeOrderInfoParam{}
	err = json.Unmarshal(request, &stationStatus)
	if err != nil {
		return nil, err
	}
	return &stationStatus, nil
}

// 请求绿喵方返回结果加密

func (c *NotifyClient) NotificationResponse(param interface{}) (resp *ChargeResponse) {
	marshal, err := json.Marshal(param)
	if err != nil {
		return &ChargeResponse{
			Ret: 500,
			Msg: "解析失败",
		}
	}
	result := c.NotificationResult(QueryResponse{
		Ret:  0,
		Msg:  nil,
		Data: marshal,
	})
	return result
}

//星星调用绿喵数据加密

func (c *NotifyClient) QueryRequestEncrypt(param interface{}) (resp *NotificationParam) {
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil
	}
	//数据加解
	pkcs5, err := encrypttool.AesEncryptPKCS5(marshal, []byte(c.AESSecret), []byte(c.AESIv))
	if err != nil {
		return
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)
	operatorID := c.OperatorID
	timestamp := fmt.Sprint(time.Now().UnixMilli()) // "20160729142400"
	seq := "0001"
	encReq := operatorID + data + timestamp + seq
	signReq := encrypttool.HMacMd5(encReq, c.SigSecret)
	queryParams := NotificationParam{
		Sig:        strings.ToUpper(signReq),
		Data:       data,
		OperatorID: operatorID,
		TimeStamp:  timestamp,
		Seq:        "0001",
	}
	return &queryParams
}

func InterfaceToString(data interface{}) string {
	var key string
	switch data.(type) {
	case string:
		key = data.(string)
	case int:
		key = strconv.Itoa(data.(int))
	case int64:
		it := data.(int64)
		key = strconv.FormatInt(it, 10)
	case float64:
		it := data.(float64)
		key = strconv.FormatFloat(it, 'f', -1, 64)
	case nil:
		key = "null"
	}
	return key
}
