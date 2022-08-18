package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"time"
)

type XingXingService struct {
	OperatorSecret string `json:"OperatorSecret,omitempty"` //运营商密钥
	OperatorID     string `json:"OperatorID,omitempty"`     //运营商标识
	SigSecret      string `json:"SigSecret,omitempty"`      //签名密钥
	DataSecret     string `json:"DataSecret,omitempty"`     //消息密钥
	DataSecretIV   string `json:"DataSecretIV,omitempty"`   //消息密钥初始化向量
	Domain         string `json:"Domain,omitempty"`         //域名
}

type getToken struct {
	OperatorSecret string `json:"OperatorSecret,omitempty"`
	OperatorID     string `json:"OperatorID,omitempty"`
	Sig            string `json:"Sig,omitempty"`
	Data           string `json:"Data,omitempty"`
	TimeStamp      string `json:"TimeStamp,omitempty"`
	Seq            string `json:"Seq,omitempty"`
}

type queryRequest struct {
	Sig        string `json:"Sig"`
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
}

// GetXingAccessToken 星星充电 query token
func (srv XingXingService) GetXingAccessToken(ctx *context.MioContext) (string, error) {
	redisCmd := app.Redis.Get(ctx, "token:"+srv.OperatorID)
	result, err := redisCmd.Result()
	if err != nil {
		return "", err
	}
	if result != "" {
		return result, nil
	}
	timeStr := time.Now().Format("20060102150405")
	//data加密
	data := getToken{
		OperatorID:     srv.OperatorID,
		OperatorSecret: srv.OperatorSecret,
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	//内容加密
	encryptData := encrypt.AesEncrypt(string(marshal), srv.DataSecret, srv.DataSecretIV)
	//签名加密
	encStr := srv.OperatorID + encryptData + timeStr + "0001"
	encryptSig := encrypt.HMacMd5(encStr, srv.SigSecret)
	queryParams := getToken{
		OperatorID: srv.OperatorID,
		Sig:        encryptSig,
		Data:       encryptData,
		TimeStamp:  timeStr,
		Seq:        "0001",
	}

	url := srv.Domain + "/query_token"
	body, err := httputil.PostJson(url, queryParams)
	if err != nil {
		return "", err
	}
	signResponse := XingSignResponse{}
	err = json.Unmarshal(body, &signResponse)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	if signResponse.Ret != 0 {
		return "", errors.New("请求错误")
	}
	//data解密
	encryptStr, _ := encrypt.AesDecrypt(signResponse.Data, srv.DataSecret, srv.DataSecretIV)
	signAccess := XingAccessResult{}
	_ = json.Unmarshal([]byte(encryptStr), &signAccess)
	//存redis
	app.Redis.Set(ctx, "token:"+srv.OperatorID, signAccess.AccessToken, time.Second*time.Duration(signAccess.TokenAvailableTime))
	return signAccess.AccessToken, nil
}

func (srv XingXingService) SendCoupon(phoneNumber string, provideId string, token string) error {
	r := struct {
		PhoneNumber string `json:"phoneNumber"`
		ProvideId   string `json:"provideId"`
	}{
		PhoneNumber: phoneNumber,
		ProvideId:   provideId,
	}
	url := srv.Domain + "/query_delivery_provide"
	contentType := httputil.HttpWithHeader("Content-Type", "application/json; charset=utf-8")
	authToken := httputil.HttpWithHeader("Authorization", token)
	body, err := httputil.PostJson(url, r, contentType, authToken)
	fmt.Printf("%s\n", body)
	if err != nil {
		return err
	}
	return nil
}
