package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"time"
)

func NewStarChargeService(context *context.MioContext) *StarChargeService {
	return &StarChargeService{
		ctx:            context,
		OperatorSecret: "3YEnj8W0negqs44Lh9ETTVEi2W1JZyt9",
		OperatorID:     "MA1FY5992", //要换
		SigSecret:      "5frdjVGMJIblh58xGNn6tQdZrBzaC9cU",
		DataSecret:     "FyTx5OwuTpEEPQJ5",
		DataSecretIV:   "ULxxy31gh7Qw67k5",
		Domain:         "https://evcs.starcharge.com/evcs/starcharge/",
		ProvideId:      "JC_20220820094600625", //要换
	}
}

type StarChargeService struct {
	ctx            *context.MioContext
	OperatorSecret string `json:"OperatorSecret,omitempty"` //运营商密钥
	OperatorID     string `json:"OperatorID,omitempty"`     //运营商标识
	SigSecret      string `json:"SigSecret,omitempty"`      //签名密钥
	DataSecret     string `json:"DataSecret,omitempty"`     //消息密钥
	DataSecretIV   string `json:"DataSecretIV,omitempty"`   //消息密钥初始化向量
	Domain         string `json:"Domain,omitempty"`         //域名
	ProvideId      string `json:"Batch,omitempty"`
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

// GetAccessToken 星星充电 query token
func (srv StarChargeService) GetAccessToken() (string, error) {
	redisCmd := app.Redis.Get(srv.ctx, "token:"+srv.OperatorID)
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
	//response
	signResponse := StarChargeResponse{}
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
	//result.data解密
	accessResult := StarChargeAccessResult{}
	encryptStr, _ := encrypt.AesDecrypt(signResponse.Data, srv.DataSecret, srv.DataSecretIV)
	_ = json.Unmarshal([]byte(encryptStr), &accessResult)
	//存redis
	app.Redis.Set(srv.ctx, "token:"+srv.OperatorID, accessResult.AccessToken, time.Second*time.Duration(accessResult.TokenAvailableTime))
	return accessResult.AccessToken, nil
}

func (srv StarChargeService) SendCoupon(openId, phoneNumber string, provideId string, token string) error {
	//保存记录
	history := entity.CouponHistory{
		OpenId:     openId,
		CouponType: "star_charge",
		CreateTime: time.Time{},
	}
	_, err := repository.DefaultCouponHistoryRepository.Insert(&history)
	if err != nil {
		return err
	}
	r := struct {
		PhoneNumber string `json:"phoneNumber"`
		ProvideId   string `json:"provideId"`
	}{
		PhoneNumber: phoneNumber,
		ProvideId:   provideId,
	}
	url := srv.Domain + "/query_delivery_provide"
	authToken := httputil.HttpWithHeader("Authorization", "Bearer "+token)
	body, err := httputil.PostJson(url, r, authToken)
	fmt.Printf("%s\n", body)
	if err != nil {
		return err
	}
	// response
	provideResponse := StarChargeResponse{}
	err = json.Unmarshal(body, &provideResponse)
	if err != nil {
		return err
	}
	if provideResponse.Ret != 0 {
		return errors.New(provideResponse.Msg)
	}
	// result.data解密
	provideResult := StarChargeProvideResult{}
	encryptStr, _ := encrypt.AesDecrypt(provideResponse.Data, srv.DataSecret, srv.DataSecretIV)
	_ = json.Unmarshal([]byte(encryptStr), &provideResult)
	if provideResult.SuccStat != 0 {
		return errors.New(provideResult.FailReasonMsg)
	}
	//更新code
	upResp := entity.CouponHistory{
		OpenId: openId,
		Code:   provideResult.CouponCode,
	}
	err = repository.DefaultCouponHistoryRepository.Update(&upResp)
	if err != nil {
		return err
	}
	return nil
}

func (srv StarChargeService) CheckLimit(openId string) error {
	builder := repository.DefaultCommentRepository.RowBuilder()
	builder.Where("open_id = ?", openId)
	_, err := repository.DefaultCommentRepository.FindOneQuery(builder)
	if err == nil {
		//已经存在
		return errors.New("每位用户只限一次")
	}
	return nil
}
