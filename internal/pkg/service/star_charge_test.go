package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"testing"
	"time"
)

func TestGetAccessToken(t *testing.T) {
	Xing := StarChargeService{
		OperatorSecret: "3YEnj8W0negqs44Lh9ETTVEi2W1JZyt9",
		OperatorID:     "MA1FY5992", //要换
		SigSecret:      "5frdjVGMJIblh58xGNn6tQdZrBzaC9cU",
		DataSecret:     "FyTx5OwuTpEEPQJ5",
		DataSecretIV:   "ULxxy31gh7Qw67k5",
		Domain:         "https://evcs.starcharge.com/evcs/starcharge/",
		ProvideId:      "JC_20220820094600625", //要换
	}

	token, err := getXingAccessToken(context.NewMioContext(), Xing)
	if err != nil {
		fmt.Printf("get token error: %e\n", err)
	}
	fmt.Printf("getXingAccessToken token: %s\n", token)
	err = sendCoupon("13083605153", Xing.ProvideId, token, Xing)
	if err != nil {
		fmt.Printf("sendCoupon error: %e\n", err)
		return
	}
	fmt.Printf("success")
}

func getXingAccessToken(ctx *context.MioContext, xing StarChargeService) (string, error) {
	data := getToken{
		OperatorSecret: xing.OperatorSecret,
		OperatorID:     xing.OperatorID,
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	//内容加密
	encryptData := encrypt.AesEncrypt(string(marshal), xing.DataSecret, xing.DataSecretIV)
	//fmt.Printf("encryptData\n%s\n-Data\n%s\n", encryptData, "j5tJ74cKFiGJ65Ot7NaSyZQoaYNUpSYy7hVWul9Yw26tXyLZb7F2Vf+58kGMk6GUfUzR6WVJn7asnFnL7UfoNg==")

	//签名加密
	timeStr := time.Now().Format("20060102150405")
	signStr := xing.OperatorID + encryptData + timeStr + "0001"
	encryptSig := encrypt.HMacMd5(signStr, xing.SigSecret)
	//fmt.Printf("encryptSig\n%s\n-Sig\n%s\n", strings.ToUpper(encryptSig), "14FC0F2D4C74CB8B4914D4038E5F4AA8")

	queryParams1 := queryRequest{
		Sig:        encryptSig,  //encryptSig,
		Data:       encryptData, //encryptData,
		OperatorID: xing.OperatorID,
		TimeStamp:  timeStr,
		Seq:        "0001",
	}

	//queryParams2 := queryRequest{
	//	Sig:        "14FC0F2D4C74CB8B4914D4038E5F4AA8",
	//	Data:       "j5tJ74cKFiGJ65Ot7NaSyZQoaYNUpSYy7hVWul9Yw26tXyLZb7F2Vf+58kGMk6GUfUzR6WVJn7asnFnL7UfoNg==",
	//	OperatorID: "MA1G55M81",
	//	TimeStamp:  "20220816153043",
	//	Seq:        "0001",
	//}
	//fmt.Printf("queryParams1\n%v\nqueryParams2\n%v\n", queryParams1, queryParams2)
	url := xing.Domain + "/query_token"
	body, err := httputil.PostJson(url, queryParams1)
	fmt.Printf("body %s\n", body)
	if err != nil {
		return "", err
	}

	signResult := StarChargeResponse{}
	err = json.Unmarshal(body, &signResult)
	if err != nil {
		return "", err
	}
	if signResult.Ret != 0 {
		return "", errors.New(signResult.Msg)
	}
	//data解密
	encryptStr, _ := encrypt.AesDecrypt(signResult.Data, xing.DataSecret, xing.DataSecretIV)
	fmt.Printf("encrypt data: %s\n", encryptStr)
	signAccess := StarChargeAccessResult{}
	_ = json.Unmarshal([]byte(encryptStr), &signAccess)
	//fmt.Printf("access response: %v\n", signAccess)
	//存redis
	//app.Redis.Set(ctx, "token:"+xing.OperatorID, signResult.AccessToken, time.Second*time.Duration(signResult.TokenAvailableTime))
	return signAccess.AccessToken, nil
}

func sendCoupon(phoneNumber string, provideId string, token string, xing StarChargeService) error {
	r := struct {
		PhoneNumber string `json:"PhoneNumber"`
		ProvideId   string `json:"ProvideId"`
	}{
		PhoneNumber: phoneNumber,
		ProvideId:   provideId,
	}
	//data加密
	marshal, _ := json.Marshal(r)
	encryptData := encrypt.AesEncrypt(string(marshal), xing.DataSecret, xing.DataSecretIV)
	//sign加密
	timeStr := time.Now().Format("20060102150405")
	signStr := xing.OperatorID + encryptData + timeStr + "0001"
	encryptSig := encrypt.HMacMd5(signStr, xing.SigSecret)

	queryParams := queryRequest{
		Sig:        encryptSig,
		Data:       encryptData,
		OperatorID: xing.OperatorID,
		TimeStamp:  timeStr,
		Seq:        "0001",
	}

	url := xing.Domain + "/query_delivery_provide"
	authToken := httputil.HttpWithHeader("Authorization", "Bearer "+token)
	body, err := httputil.PostJson(url, queryParams, authToken)
	fmt.Printf("sendcoupon body %s\n", body)
	if err != nil {
		return err
	}
	//data解密
	signResult := StarChargeResponse{}
	err = json.Unmarshal(body, &signResult)
	if err != nil {
		return err
	}
	if signResult.Ret != 0 {
		return errors.New(signResult.Msg)
	}

	encryptStr, _ := encrypt.AesDecrypt(signResult.Data, xing.DataSecret, xing.DataSecretIV)
	fmt.Printf("sendcoupon encrypt data: %s\n", encryptStr)
	//最终解密
	return nil
}
