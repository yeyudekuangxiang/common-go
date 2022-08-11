package initialize

import (
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/pkg/wxoa"
)

func InitWxoa() {
	client := wxoa.NewWxOA(config.Config.MioSrvOA.AppId, config.Config.MioSrvOA.Secret, app.Redis)
	*app.WxOa = *client
}
