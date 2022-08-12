package activity

import (
	"encoding/json"
	"errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/util/httputil"
	"time"
)

//仅用于获取第三方平台认证

type AuthService struct {
	XingAccessReq XingAccessRequest
}

//星星充电 query token
func (srv AuthService) GetXingAccessToken(ctx *context.MioContext, url string) (string, error) {
	redisCmd := app.Redis.Get(ctx, "token:"+srv.XingAccessReq.OperatorID)
	result, err := redisCmd.Result()
	if err != nil {
		return "", err
	}
	if result != "" {
		return result, nil
	}
	body, err := httputil.PostJson(url, srv.XingAccessReq)
	if err != nil {
		return "", err
	}
	res := XingAccess{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if res.FailReason != 0 {
		return "", errors.New("请求错误")
	}
	//存redis
	app.Redis.Set(ctx, "token:"+srv.XingAccessReq.OperatorID, res.AccessToken, time.Second*time.Duration(res.TokenAvailableTim))
	return res.AccessToken, nil
}
