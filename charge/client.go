package charge

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"log"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Domain     string
	Version    string
	AESSecret  string //acb93539fc9bg78k
	AESIv      string //82c91325e74bef0f
	SigSecret  string //9af2e7b2d7562ad5  签名密钥
	Token      string
	OperatorID string

	MIOAESSecret string //绿喵AESSecret
	MIOAESIv     string //绿喵AESIv
	MIOSigSecret string //9af2e7b2d7562ad5  签名密钥
}

//调用充电接口

func (c *Client) Request(param SendChargeParam) (resp *QueryResponse, err error) {
	sendUrl := fmt.Sprintf("%s%s", c.Domain, param.QueryUrl)
	//数据加解
	pkcs5, err := encrypttool.AesEncryptPKCS5(param.Data, []byte(c.AESSecret), []byte(c.AESIv))
	if err != nil {
		return
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)
	operatorID := c.OperatorID
	timestamp := fmt.Sprint(time.Now().UnixMilli()) // "20160729142400"
	seq := "0001"
	encReq := operatorID + data + timestamp + seq
	signReq := encrypttool.HMacMd5(encReq, c.SigSecret)
	queryParams := QueryRequest{
		Sig:        strings.ToUpper(signReq),
		Data:       data,
		OperatorID: operatorID,
		TimeStamp:  timestamp,
		Seq:        "0001",
	}
	queryParamsstr, err := json.Marshal(queryParams)
	println(queryParamsstr)
	authToken := httptool.HttpWithHeader("Authorization", "Bearer "+c.Token)
	body, err := httptool.PostJson(sendUrl, queryParams, authToken)
	if err != nil {
		log.Printf("request:%s response:%v\n", queryParams, err)
		return nil, err
	}
	log.Printf("request:%s response:%s\n", queryParams, string(body))

	res := &ChargeResponse{}
	err = json.Unmarshal(body, &res)

	if res.Ret != 0 {
		return nil, errors.New(c.interfaceToString(res.Msg))
	}
	encResp := strconv.Itoa(res.Ret) + c.interfaceToString(res.Msg) + res.Data
	signResp := encrypttool.HMacMd5(encResp, c.SigSecret)
	if strings.ToUpper(signResp) != res.Sig {
		return nil, errors.New("签名有误")
	}
	decodeString, err := base64.StdEncoding.DecodeString(res.Data)
	if err != nil {
		return nil, err
	}
	decodeData, err := encrypttool.AesDecryptPKCS5(decodeString, []byte(c.AESSecret), []byte(c.AESIv))
	if err != nil {
		return nil, err
	}
	resV2 := &QueryResponse{
		Ret:  res.Ret,
		Msg:  res.Msg,
		Data: decodeData,
		Sig:  res.Sig,
	}
	return resV2, nil
}

//请求星星接口

func (c *Client) QueryToken(param QueryTokenParam) (resp *QueryTokenResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_token",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryTokenResult{}
	err = json.Unmarshal(response.Data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//请求设备认证

func (c *Client) QueryEquipAuth(param QueryEquipAuthParam) (resp *QueryEquipAuthResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_equip_auth",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryEquipAuthResult{}
	err = json.Unmarshal(response.Data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil

}

//请求启动充电

func (c *Client) QueryStartCharge(param QueryStartChargeParam) (resp *QueryStartChargeResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_start_charge",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryStartChargeResult{}
	err = json.Unmarshal(response.Data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//查询充电状态

func (c *Client) QueryEquipChargeStatus(param QueryEquipChargeStatusParam) (resp *QueryEquipChargeStatusResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_equip_charge_status",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryEquipChargeStatusResult{}
	err = json.Unmarshal(response.Data, ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//请求停止充电

func (c *Client) QueryStopCharge(param QueryStopChargeParam) (resp *QueryStopChargeResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_stop_charge",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryStopChargeResult{}
	err = json.Unmarshal(response.Data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//获取充电信息

func (c *Client) QueryStationsInfo(param QueryStationsInfoParam) (resp *QueryStationsInfoResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_stations_info",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryStationsInfoResult{}
	err = json.Unmarshal(response.Data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//设备接口状态查询

func (c *Client) QueryStationStatus(param QueryStationStatusParam) (resp *QueryStationStatusResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendChargeParam{
		data,
		"query_station_status",
	})
	if err != nil {
		return nil, err
	}
	ret := QueryStationStatusResult{}
	err = json.Unmarshal(response.Data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil

}

//星星回调绿喵接口

func (c *Client) NotificationRequest(param NotificationParam) (resp []byte, err error) {
	operatorID := param.OperatorID
	data := param.Data
	timestamp := param.TimeStamp
	seq := param.Seq

	encReq := operatorID + data + timestamp + seq
	sign := strings.ToUpper(encrypttool.HMacMd5(encReq, c.MIOSigSecret))

	if sign != param.Sig {
		return nil, errors.New("签名失败")
	}
	dataString, err := base64.StdEncoding.DecodeString(param.Data)
	if err != nil {
		return nil, err
	}
	//数据解密
	decryptPKCS5, err := encrypttool.AesDecryptPKCS5(dataString, []byte(c.MIOAESSecret), []byte(c.MIOAESIv))
	if err != nil {
		return nil, err
	}
	return decryptPKCS5, nil
}

//返回结果加密返回

func (c *Client) NotificationResult(param QueryResponse) (resp *ChargeResponse) {
	if param.Ret != 0 {
		return &ChargeResponse{
			Ret:  param.Ret,
			Msg:  param.Msg,
			Data: "",
			Sig:  "",
		}
	}
	pkcs5, err := encrypttool.AesEncryptPKCS5(param.Data, []byte(c.MIOAESSecret), []byte(c.MIOAESIv))
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
	encResp := strconv.Itoa(param.Ret) + c.interfaceToString(param.Msg) + data
	sign := encrypttool.HMacMd5(encResp, c.MIOSigSecret)
	res := ChargeResponse{
		Ret:  param.Ret,
		Msg:  param.Msg,
		Data: data,
		Sig:  sign,
	}
	return &res
}

//请求绿喵token

func (c *Client) QueryTokenRequest(param NotificationParam) (resp *QueryTokenReq, err error) {
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

//请求绿喵token返回结果

func (c *Client) QueryTokenResult(param QueryTokenResp) (resp *ChargeResponse, err error) {
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	result := c.NotificationResult(QueryResponse{
		Ret:  0,
		Msg:  nil,
		Data: marshal,
		Sig:  "",
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

//设备状态变化推送

func (c *Client) NotificationStationStatusRequest(param NotificationParam) (resp *NotificationStationStatusParam, err error) {
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

func (c *Client) NotificationStartChargeResultRequest(param NotificationParam) (resp *NotificationStartChargeResultParam, err error) {
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

func (c *Client) NotificationEquipChargeStatusRequest(param NotificationParam) (resp *NotificationEquipChargeStatusParam, err error) {
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

func (c *Client) NotificationStopChargeResultRequest(param NotificationParam) (resp *NotificationStopChargeResultParam, err error) {
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

func (c *Client) NotificationChargeOrderInfoRequest(param NotificationParam) (resp *NotificationChargeOrderInfoParam, err error) {
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

func (c *Client) NotificationResponse(param interface{}) (resp *ChargeResponse) {
	marshal, err := json.Marshal(param)
	if err != nil {
		return &ChargeResponse{
			Ret:  500,
			Msg:  "解析失败",
			Data: "",
			Sig:  "",
		}
	}
	result := c.NotificationResult(QueryResponse{
		Ret:  0,
		Msg:  nil,
		Data: marshal,
		Sig:  "",
	})
	return result
}

func (c *Client) interfaceToString(data interface{}) string {
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
