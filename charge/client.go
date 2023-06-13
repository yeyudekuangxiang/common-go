package charge

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"log"
	"math/rand"
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
}

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
		return nil, errors.New(InterfaceToString(res.Msg))
	}
	encResp := strconv.Itoa(res.Ret) + InterfaceToString(res.Msg) + res.Data
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

//获取流水号订单号

func GenerateSerialNumber(operatorID string) string {
	rand.Seed(time.Now().UnixNano())
	uniqueID := rand.Intn(1000000000)
	serialNumber := fmt.Sprintf("%s%09d", operatorID, uniqueID)
	return serialNumber
}

//获取StartChargeSeqStat对应的中文状态

func GetStartChargeSeqStatDesc(startChargeSeqStat int) string {
	switch startChargeSeqStat {
	case 1:
		return "启动中"
	case 2:
		return "充电中"
	case 3:
		return "停止中"
	case 4:
		return "已结束"
	case 5:
		return "未知"
	default:
		return "未知"
	}
}

//获取FailReason对应的中文状态

func GetFailReasonDesc(FailReason int) string {
	switch FailReason {
	case 0:
		return "无"
	case 1:
		return "此设备尚未插抢"
	case 2:
		return "设备检测失败"
	default:
		return "未知错误"
	}
}
