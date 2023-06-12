package starcharge

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"log"
	"time"
)

const (
	//每个接口的URL均采用如下格式定义：http（s）： //【域名】/evcs/v【版本号】/【接口名称】a） 域名：各接入运营商所属域名。
	version        = "1.0"
	actionBikeCard = "hellobike.activity.bikecard"
	url            = "https://%s/evcs/%s/%s"
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

// SendCoupon 发放电子票

type queryRequest struct {
	Sig        string `json:"Sig"`
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
}

//请求设备认证

func (c *Client) QueryEquipAuth(param QueryEquipAuthParam) (resp *QueryEquipAuthResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendStarChargeParam{
		data,
		"query_equip_auth",
	})
	if err != nil {
		return nil, err
	}
	if response.Ret == 0 {
		return nil, errors.New(response.Msg)
	}
	ret := QueryEquipAuthResult{}
	err = json.Unmarshal(response.Data, ret)
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
	response, err := c.Request(SendStarChargeParam{
		data,
		"query_start_charge",
	})
	if err != nil {
		return nil, err
	}
	if response.Ret == 0 {
		return nil, errors.New(response.Msg)
	}
	ret := QueryStartChargeResult{}
	err = json.Unmarshal(response.Data, ret)
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
	response, err := c.Request(SendStarChargeParam{
		data,
		"query_equip_charge_status",
	})
	if err != nil {
		return nil, err
	}
	if response.Ret == 0 {
		return nil, errors.New(response.Msg)
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
	response, err := c.Request(SendStarChargeParam{
		data,
		"query_stop_charge",
	})
	if err != nil {
		return nil, err
	}
	if response.Ret == 0 {
		return nil, errors.New(response.Msg)
	}
	ret := QueryStopChargeResult{}
	err = json.Unmarshal(response.Data, ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) QueryStationsInfo(param QueryStationsInfoParam) (resp *QueryStationsInfoResult, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	response, err := c.Request(SendStarChargeParam{
		data,
		"query_stations_info",
	})
	if err != nil {
		return nil, err
	}
	if response.Ret == 0 {
		return nil, errors.New(response.Msg)
	}
	ret := QueryStationsInfoResult{}
	err = json.Unmarshal(response.Data, ret)
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
	response, err := c.Request(SendStarChargeParam{
		data,
		"query_station_status",
	})
	if err != nil {
		return nil, err
	}
	if response.Ret == 0 {
		return nil, errors.New(response.Msg)
	}
	ret := QueryStationStatusResult{}
	err = json.Unmarshal(response.Data, ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil

}

//请求星星接口

func (c *Client) Request(param SendStarChargeParam) (resp *starChargeResponse, err error) {
	sendUrl := fmt.Sprintf(url, c.Domain, c.Version, param.QueryUrl)
	//数据加解
	pkcs5, err := encrypttool.AesEncryptPKCS5(param.Data, []byte(c.AESSecret), []byte(c.AESIv))
	if err != nil {
		return
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)

	//签名
	operatorID := c.OperatorID
	// 示例参数信息（Data）
	data = "il7B0BSEjFdzpyKzfOFpvg/Se1CP802RItKYFPfSLRxJ3jfObV19hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlqsJauko79NnwQJbzDTyLooYolwz75qBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9tNqs/e2Vjja/ltE1M0lqvxfXQ6da6HrThsm5id4ClZFliOacRfrsPLRixS/IQYtksxghvJwbqOsblsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk713wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	// 示例时间戳（TimeStamp）
	timestamp := fmt.Sprint(time.Now().UnixMilli()) // "20160729142400"
	// 示例自增序列（Seq）
	seq := "0001"

	sign := getHMACMD5Signature(operatorID, data, timestamp, seq, c.SigSecret)
	queryParams := queryRequest{
		Sig:        sign,
		Data:       data,
		OperatorID: operatorID,
		TimeStamp:  timestamp,
		Seq:        "0001",
	}

	authToken := httptool.HttpWithHeader("Authorization", "Bearer "+c.Token)
	body, err := httptool.PostJson(sendUrl, queryParams, authToken)
	if err != nil {
		log.Printf("request:%s response:%v\n", queryParams, err)
		return nil, err
	}
	log.Printf("request:%s response:%s\n", queryParams, string(body))
	res := &starChargeResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//星星回调绿喵接口

func (c *Client) NotificationValidate(param NotificationParam) (resp *[]byte, err error) {
	operatorID := param.OperatorID // 示例参数信息（Data）
	data := param.Data
	timestamp := param.TimeStamp // 示例时间戳（TimeStamp）
	seq := param.Seq             // 示例自增序列（Seq）
	sign := getHMACMD5Signature(operatorID, data, timestamp, seq, c.MIOSigSecret)
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
	return &decryptPKCS5, nil
}

// HMAC-MD5 参数签名
func getHMACMD5Signature(operatorID string, data string, timeStamp string, seq string, sigSecret string) string {
	// 拼接参数
	message := operatorID + data + timeStamp + seq
	// 计算签名
	key := []byte(sigSecret)
	payload := []byte(message)
	hash := hmac.New(md5.New, key)
	hash.Write(payload)
	signature := hash.Sum(nil)
	// 将签名转换为大写字符串
	return hex.EncodeToString(signature)
}
