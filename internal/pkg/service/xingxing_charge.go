package service

import (
	"encoding/json"
	"errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	"time"
)

type XingXingService struct {
	OperatorID     string `json:"operatorID,omitempty"`
	OperatorSecret string `json:"operatorSecret,omitempty"`
	SigSecret      string `json:"sigSecret,omitempty"`
	DataSecret     string `json:"dataSecret,omitempty"`
	DataSecretIV   string `json:"dataSecretIV,omitempty"`
	Url            string `json:"url"`
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
	r := struct {
		OperatorID     string `json:"operatorID"`
		OperatorSecret string `json:"operatorSecret"`
	}{
		OperatorID:     srv.OperatorID,
		OperatorSecret: srv.OperatorSecret,
	}
	url := srv.Url + "/query_token"
	body, err := httputil.PostJson(url, r)
	if err != nil {
		return "", err
	}
	signResult := XingXingSignResult{}
	err = json.Unmarshal(body, &signResult)
	if err != nil {
		return "", err
	}
	if signResult.FailReason != 0 {
		return "", errors.New("请求错误")
	}
	//存redis
	app.Redis.Set(ctx, "token:"+srv.OperatorID, signResult.AccessToken, time.Second*time.Duration(signResult.TokenAvailableTime))
	return signResult.AccessToken, nil
}

func (srv XingXingService) SendCount(ctx *context.MioContext, phoneNumber string, provideId string, token string) error {
	r := struct {
		PhoneNumber string `json:"phoneNumber"`
		ProvideId   string `json:"provideId"`
	}{
		PhoneNumber: phoneNumber,
		ProvideId:   provideId,
	}
	url := srv.Url + "/query_delivery_provide"
	contentType := httputil.HttpWithHeader("Content-Type", "application/json; charset=utf-8")
	authToken := httputil.HttpWithHeader("Authorization", token)
	body, err := httputil.PostJson(url, r, contentType, authToken)
	if err != nil {
		return err
	}
	result := XingXingSendCouponResult{}
	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}
	if result.SuccStat == 1 {
		return errno.ErrAuth
	}
	return nil
}
