package wxapp

import (
	"context"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"testing"
	"time"
)

func TestAutoTry(t *testing.T) {

	var qrResp *QRCodeResponse

	client := NewClient("wx3279a0b1782e2d93", "4028279bec13ece9155c7eae348de171", NewTokenCenter())

	err := client.AutoTryAccessToken(func(accessToken string) (try bool, err error) {
		qrResp, err = client.GetUnlimitedQRCodeResponse(&weapp.UnlimitedQRCode{
			Scene: "123",
			Page:  "pages/community/details/index",
			Width: 100,
		})

		//系统错误
		if err != nil {
			return false, err
		}

		//自动判断是否重试
		isExpire, err := client.IsExpireAccessToken(qrResp.ErrCode)
		return isExpire, err
	}, 3)

	assert.Equal(t, nil, err)

	fmt.Println("请求结果", qrResp.ErrCode, qrResp.ErrMsg)
}

type TokenCenter struct {
	rpc tokencenterclient.TokenCenter
}

func NewTokenCenter() *TokenCenter {
	return &TokenCenter{
		rpc: tokencenterclient.NewTokenCenter(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: []string{"127.0.0.1:1000"},
			NonBlock:  true,
		})),
	}
}

func (t TokenCenter) AccessToken() (token string, expireIn time.Time, err error) {
	tokenResp, err := t.rpc.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: "1",
		OldToken: "",
		Refresh:  false,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenResp.AccessToken, time.UnixMilli(tokenResp.ExpireAt), nil
}

func (t TokenCenter) RefreshToken(oldToken string) (token string, expireIn time.Time, err error) {
	tokenResp, err := t.rpc.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: "1",
		OldToken: oldToken,
		Refresh:  true,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenResp.AccessToken, time.UnixMilli(tokenResp.ExpireAt), nil
}

func (t TokenCenter) IsExpired(code string) (bool, error) {
	tokenResp, err := t.rpc.IsAccessTokenExpired(context.Background(), &tokencenterclient.IsAccessTokenExpiredReq{
		Code:     code,
		CenterId: "1",
	})
	if err != nil {
		return false, err
	}
	return tokenResp.IsExpired, nil
}
