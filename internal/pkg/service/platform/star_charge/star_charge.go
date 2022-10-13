package star_charge

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
		OperatorID:     "MA1FY5992",
		SigSecret:      "5frdjVGMJIblh58xGNn6tQdZrBzaC9cU",
		DataSecret:     "FyTx5OwuTpEEPQJ5",
		DataSecretIV:   "ULxxy31gh7Qw67k5",
		Domain:         "https://evcs.starcharge.com/evcs/starcharge/",
		ProvideId:      "JC_20220820094600625",
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
	token := app.Redis.Get(srv.ctx, "token:"+srv.OperatorID).Val()
	if token != "" {
		return token, nil
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
	signResponse := starChargeResponse{}
	err = json.Unmarshal(body, &signResponse)
	if err != nil {
		return "", err
	}
	if signResponse.Ret != 0 {
		return "", errors.New("请求错误")
	}
	//result.data解密
	accessResult := starChargeAccessResult{}
	encryptStr, _ := encrypt.AesDecrypt(signResponse.Data, srv.DataSecret, srv.DataSecretIV)
	_ = json.Unmarshal([]byte(encryptStr), &accessResult)
	//存redis
	app.Redis.Set(srv.ctx, "token:"+srv.OperatorID, accessResult.AccessToken, time.Second*time.Duration(accessResult.TokenAvailableTime))
	return accessResult.AccessToken, nil
}

func (srv StarChargeService) SendCoupon(openId, phoneNumber string, provideId string, token string) error {
	r := struct {
		PhoneNumber string `json:"PhoneNumber"`
		ProvideId   string `json:"ProvideId"`
	}{
		PhoneNumber: phoneNumber,
		ProvideId:   provideId,
	}
	//data加密
	marshal, _ := json.Marshal(r)
	encryptData := encrypt.AesEncrypt(string(marshal), srv.DataSecret, srv.DataSecretIV)
	//sign加密
	timeStr := time.Now().Format("20060102150405")
	signStr := srv.OperatorID + encryptData + timeStr + "0001"
	encryptSig := encrypt.HMacMd5(signStr, srv.SigSecret)
	queryParams := queryRequest{
		Sig:        encryptSig,
		Data:       encryptData,
		OperatorID: srv.OperatorID,
		TimeStamp:  timeStr,
		Seq:        "0001",
	}
	url := srv.Domain + "/query_delivery_provide"
	authToken := httputil.HttpWithHeader("Authorization", "Bearer "+token)
	body, err := httputil.PostJson(url, queryParams, authToken)
	fmt.Printf("%s\n", body)
	if err != nil {
		return err
	}
	// response
	provideResponse := starChargeResponse{}
	err = json.Unmarshal(body, &provideResponse)
	if err != nil {
		return err
	}
	if provideResponse.Ret != 0 {
		return errors.New(provideResponse.Msg)
	}
	// result.data解密
	provideResult := starChargeProvideResult{}
	encryptStr, _ := encrypt.AesDecrypt(provideResponse.Data, srv.DataSecret, srv.DataSecretIV)
	_ = json.Unmarshal([]byte(encryptStr), &provideResult)
	if provideResult.SuccStat != 0 {
		return errors.New(provideResult.FailReasonMsg)
	}
	//保存记录
	history := entity.CouponHistory{
		OpenId:     openId,
		CouponType: "star_charge",
		Code:       provideResult.CouponCode,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	_, err = repository.DefaultCouponHistoryRepository.Insert(&history)
	if err != nil {
		fmt.Printf("星星充电,insert error:%s", err.Error())
		return err
	}
	return nil
}

// CheckChargeLimit 充电检测
func (srv StarChargeService) CheckChargeLimit(openId string, startTime, endTime string) error {
	todayBuilder := repository.DefaultCouponHistoryRepository.RowBuilder()
	todayBuilder.Where("open_id = ?", openId).
		Where("coupon_type = ?", "star_charge").
		Where("date(create_time) = CURRENT_DATE")
	count, err := repository.DefaultCouponHistoryRepository.FindCount(todayBuilder)
	if err != nil {
		return err
	}
	if count >= 1 {
		return errors.New("每日每位用户限制领取 1 次")
	}
	builder := repository.DefaultCouponHistoryRepository.RowBuilder()
	builder.Where("open_id = ?", openId).
		Where("coupon_type = ?", "star_charge").
		Where("create_time > ?", startTime).
		Where("create_time < ?", endTime)
	count, err = repository.DefaultCouponHistoryRepository.FindCount(builder)
	if err != nil {
		return err
	}
	if count >= 2 {
		return errors.New("活动期间每位用户限制领取 2 次")
	}
	return nil

}
