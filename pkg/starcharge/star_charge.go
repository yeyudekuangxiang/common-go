package starcharge

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"time"
)

type StarCharge struct {
	ctx          *context.MioContext
	timeStr      string
	redis        *redis.Client
	dataSecret   string
	dataSecretIV string
	domain       string
	provideId    string
	sigSecret    string
	getParams
}

func NewStarCharge(ctx *context.MioContext, operatorSecret, operatorID, sigSecret, dataSecret, dataSecretIV, provideId, domain string, client *redis.Client) *StarCharge {
	return &StarCharge{
		ctx:          ctx,
		timeStr:      time.Now().Format("20060102150405"),
		redis:        client,
		dataSecret:   dataSecret,
		dataSecretIV: dataSecretIV,
		domain:       domain,
		provideId:    provideId,
		sigSecret:    sigSecret,
		getParams: getParams{
			OperatorSecret: operatorSecret,
			OperatorID:     operatorID,
		},
	}
}

type getParams struct {
	OperatorSecret string `json:"OperatorSecret,omitempty"`
	OperatorID     string `json:"OperatorID,omitempty"`
}

type getTokenRequest struct {
	getParams
	Sig       string `json:"Sig,omitempty"`
	Data      string `json:"Data,omitempty"`
	TimeStamp string `json:"TimeStamp,omitempty"`
	Seq       string `json:"Seq,omitempty"`
}

type starChargeResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data string `json:"Data"`
	Sig  string `json:"Sig"`
}

type starChargeGetTokenResponse struct {
	OperatorID         string `json:"operatorID,omitempty"`
	SucStat            int    `json:"sucStat,omitempty"`
	AccessToken        string `json:"accessToken,omitempty"`
	TokenAvailableTime int    `json:"tokenAvailableTime,omitempty"`
	FailReason         int    `json:"failReason,omitempty"`
}

func (s StarCharge) GetAccessToken() (string, error) {
	token := s.redis.Get(s.ctx, "star_charge_token:"+s.OperatorID).Val()
	if token != "" {
		return token, nil
	}
	encryptData, _ := s.encodeData()
	encodeSign := s.encodeSign(encryptData, s.timeStr)
	//request 参数
	getTokenRequestParams := getTokenRequest{
		getParams: getParams{
			OperatorID: s.OperatorID,
		},
		Sig:       encodeSign,
		Data:      encryptData,
		TimeStamp: s.timeStr,
		Seq:       "0001",
	}
	url := s.domain + "/query_token"
	body, err := httputil.PostJson(url, getTokenRequestParams)
	if err != nil {
		return "", err
	}
	//decode response
	response, err := s.decodeResponse(body)
	if err != nil {
		return "", err
	}
	//decode response.data
	result, err := s.decodeData(response.Data)
	if err != nil {
		return "", err
	}
	//存redis
	s.redis.Set(s.ctx, "star_charge_token:"+s.OperatorID, result.AccessToken, time.Second*time.Duration(result.TokenAvailableTime))
	return result.AccessToken, nil
}

func (s StarCharge) decodeResponse(body []byte) (starChargeResponse, error) {
	signResponse := starChargeResponse{}
	if err := json.Unmarshal(body, &signResponse); err != nil {
		return signResponse, err
	}
	if signResponse.Ret != 0 {
		return signResponse, errors.New("请求错误")
	}
	return signResponse, nil
}

func (s StarCharge) decodeData(data string) (starChargeGetTokenResponse, error) {
	accessResult := starChargeGetTokenResponse{}
	encryptStr, _ := encrypt.AesDecrypt(data, s.dataSecret, s.dataSecretIV)
	_ = json.Unmarshal([]byte(encryptStr), &accessResult)
	if accessResult.SucStat == 1 {
		return accessResult, errors.New("获取token错误")
	}
	return accessResult, nil
}

func (s StarCharge) encodeData() (string, error) {
	//data加密
	marshal, err := json.Marshal(s.getParams)
	if err != nil {
		return "", err
	}
	//内容加密
	return encrypt.AesEncrypt(string(marshal), s.dataSecret, s.dataSecretIV), nil
}

func (s StarCharge) encodeSign(encryptData, timeStr string) string {
	//签名加密
	encStr := s.OperatorID + encryptData + timeStr + "0001"
	return encrypt.HMacMd5(encStr, s.sigSecret)
}
