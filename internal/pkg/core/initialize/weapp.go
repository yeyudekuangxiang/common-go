package initialize

import (
	"context"
	"github.com/medivhzhan/weapp/v3/logger"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/wxapp"
	"time"
)

//debug info warn error
var logLevelMap = map[string]logger.Level{
	"debug": logger.Info,
	"info":  logger.Info,
	"warn":  logger.Warn,
	"error": logger.Error,
}

func InitWeapp() {
	log.Println("初始化weapp组件...")
	weappSetting := config.Config.Weapp
	client := wxapp.NewClient(weappSetting.AppId, weappSetting.Secret, NewTokenCenter(), logLevelMap[config.Config.Log.Level])
	*app.Weapp = *client
	log.Println("初始化weapp组件成功")
}

type TokenCenter struct {
}

func NewTokenCenter() *TokenCenter {
	return &TokenCenter{}
}

func (t TokenCenter) AccessToken() (token string, expireIn time.Time, err error) {
	tokenResp, err := app.RpcService.TokenCenterRpcSrv.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
		CenterId: "1",
		OldToken: "",
		Refresh:  false,
	})
	if err != nil {
		return "", time.Time{}, err
	}

	return "62_LyFbC3E1ZjFL-20YUBdE4vbMj_3TXPEijjDaOi5nfj-JTpkuw2YlsshIYb8NE-CU72bu4VAR0x4VT4XAe9u1h_B4ZxVtOLeWHLJemRHrf8pJW8RtkeigGw3gEU-z08CJcHCKiw7ypf4qLC-_RIOhAAAIZF", time.UnixMilli(tokenResp.ExpireAt), nil
}

func (t TokenCenter) RefreshToken(oldToken string) (token string, expireIn time.Time, err error) {
	tokenResp, err := app.RpcService.TokenCenterRpcSrv.AccessToken(context.Background(), &tokencenterclient.GetAccessTokenReq{
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
	tokenResp, err := app.RpcService.TokenCenterRpcSrv.IsAccessTokenExpired(context.Background(), &tokencenterclient.IsAccessTokenExpiredReq{
		Code:     code,
		CenterId: "1",
	})
	if err != nil {
		return false, err
	}
	return tokenResp.IsExpired, nil
}
