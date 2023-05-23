package starcharge

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/encrypttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"log"
	"net/http"
	"time"
)

const (
	//每个接口的URL均采用如下格式定义：http（s）： //【域名】/evcs/v【版本号】/【接口名称】a） 域名：各接入运营商所属域名。
	version        = "1.0"
	actionBikeCard = "hellobike.activity.bikecard"
	url            = "https://%s/evcs/%s/%s"
)

type Client struct {
	htpClient  http.Client
	Domain     string
	Version    string
	AESSecret  string
	SigSecret  string
	Token      string
	OperatorID string
}

// SendCoupon 发放电子票

type queryRequest struct {
	Sig        string `json:"Sig"`
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
}

func (c *Client) Request(param SendStarChargeParam) (resp *starChargeResponse, err error) {
	sendUrl := fmt.Sprintf(url, c.Domain, c.Version, param.QueryUrl)
	//数据加解
	pkcs5, err := encrypttool.AesEncryptPKCS5([]byte(param.Data), []byte(c.AESSecret))
	if err != nil {
		return
	}
	data := base64.StdEncoding.EncodeToString(pkcs5)

	//签名
	operatorID := c.OperatorID
	// 示例参数信息（Data）
	data = "il7B0BSEjFdzpyKzfOFpvg/Se1CP802RItKYFPfSLRxJ3jfObV19hvYOEktPAYW2nd7S8MBcyHYyacHKbISq5iTmDzG+ivnR+SZJv3USNTYVMz9rCQVSxd0cLlqsJauko79NnwQJbzDTyLooYolwz75qBOH2/xOMirpeEqRJrF/EQjWekJmGk9RtboXePu2rka+Xm51syBPhiXJAq0GfbfaFu9tNqs/e2Vjja/ltE1M0lqvxfXQ6da6HrThsm5id4ClZFliOacRfrsPLRixS/IQYtksxghvJwbqOsblsITail9Ayy4tKcogeEZiOO+4Ed264NSKmk713wKwJLAFjCFogBx8GE3OBz4pqcAn/ydA="
	// 示例时间戳（TimeStamp）
	timestamp := string(time.Now().UnixMilli()) // "20160729142400"
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
