package star_charge

import (
	"encoding/json"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util/encrypt"
	"mio/internal/pkg/util/httputil"
	"mio/internal/pkg/util/limit"
	"mio/pkg/errno"
	"time"
)

func NewStarChargeService(ctx *context.MioContext) *StarChargeService {
	return &StarChargeService{
		ctx:            ctx,
		history:        repository.NewCouponHistoryRepository(ctx),
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
	history        repository.CouponHistoryModel
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
		return "", errno.ErrCommon.WithMessage("请求错误")
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
		return errno.ErrCommon.WithMessage(provideResponse.Msg)
	}
	// result.data解密
	provideResult := starChargeProvideResult{}
	encryptStr, _ := encrypt.AesDecrypt(provideResponse.Data, srv.DataSecret, srv.DataSecretIV)
	_ = json.Unmarshal([]byte(encryptStr), &provideResult)
	if provideResult.SuccStat != 0 {
		return errno.ErrCommon.WithMessage(provideResult.FailReasonMsg)
	}

	//保存记录
	_, err = srv.history.Insert(&entity.CouponHistory{
		OpenId:     openId,
		CouponType: "star_charge",
		Code:       provideResult.CouponCode,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})

	if err != nil {
		app.Logger.Errorf("星星充电发券记录插入失败:%s", err.Error())
	}

	return nil
}

// CheckChargeLimit 充电检测
func (srv StarChargeService) CheckChargeLimit(openId string, endTime time.Time) error {
	keyPrefix := "periodLimit:sendCoupon:star_charge:" + time.Now().Format("20060102") + ":"
	periodLimit := limit.NewPeriodLimit(int(time.Hour.Seconds())*24, 1, app.Redis, keyPrefix, limit.Align())
	res1, err := periodLimit.TakeCtx(srv.ctx.Context, openId)
	if err != nil {
		return err
	}

	if res1 != 1 {
		return errno.ErrCommon.WithMessage("每日上限1次")
	}

	keyPrefix2 := "periodLimit:sendCoupon:star_charge:"
	periodLimit2 := limit.NewPeriodLimit(int(endTime.Sub(time.Now()).Seconds()), 2, app.Redis, keyPrefix2, limit.Align())
	res2, err := periodLimit2.TakeCtx(srv.ctx.Context, openId)
	if err != nil {
		return err
	}

	if res2 != 1 {
		return errno.ErrCommon.WithMessage("活动内上限2次")
	}

	return nil
}
