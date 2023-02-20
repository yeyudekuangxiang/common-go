package initialize

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/wxoa"
	"mio/config"
	"mio/internal/pkg/core/app"
)

func InitWxoa() {
	client := wxoa.NewWxOA(config.Config.MioSrvOA.AppId, config.Config.MioSrvOA.Secret, config.RedisKey.AccessToken, app.Redis)
	*app.WxOa = *client
}
