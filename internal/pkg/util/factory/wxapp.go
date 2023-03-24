package factory

import (
	"context"
	"fmt"
	"github.com/medivhzhan/weapp/v3/logger"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"time"
)
import "gitlab.miotech.com/miotech-application/backend/common-go/wxapp"

func NewWxAppFromTokenCenterRpc(centerId string, centerRpc tokencenterclient.TokenCenter, logLevel logger.Level) (*wxapp.Client, error) {
	infoResp, err := centerRpc.GetCenterConfig(context.Background(), &tokencenterclient.GetCenterConfigReq{
		CenterId: centerId,
	})
	if err != nil {
		return nil, err
	}
	if !infoResp.Exist {
		return nil, fmt.Errorf("centerId%s不存在", centerId)
	}

	return wxapp.NewClient(infoResp.Info.AccessKey, infoResp.Info.AccessSecret, &weAppTokenCenter{
		centerID:       centerId,
		tokenCenterRpc: centerRpc,
	}, logLevel), nil
}

type weAppTokenCenter struct {
	centerID       string
	tokenCenterRpc tokencenterclient.TokenCenter
}

func (t weAppTokenCenter) AccessToken() (token string, expireIn time.Time, err error) {
	tokenResp, err := t.tokenCenterRpc.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: t.centerID,
		OldToken: "",
		Refresh:  false,
	})
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenResp.AccessToken, time.UnixMilli(tokenResp.ExpireAt), nil
}

func (t weAppTokenCenter) RefreshToken(oldToken string) (token string, expireIn time.Time, err error) {
	tokenResp, err := t.tokenCenterRpc.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: t.centerID,
		OldToken: oldToken,
		Refresh:  true,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenResp.AccessToken, time.UnixMilli(tokenResp.ExpireAt), nil
}

func (t weAppTokenCenter) IsExpired(code string) (bool, error) {
	tokenResp, err := t.tokenCenterRpc.IsAccessTokenExpired(context.Background(), &tokencenterclient.IsAccessTokenExpiredReq{
		Code:     code,
		CenterId: t.centerID,
	})
	if err != nil {
		return false, err
	}
	return tokenResp.IsExpired, nil
}
