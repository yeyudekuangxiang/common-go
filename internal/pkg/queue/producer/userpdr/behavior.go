package userpdr

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/producer/pdr"
	"mio/internal/pkg/queue/types/message/usermsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func Interaction(rk routerkey.BehaviorRouterKey, info usermsg.Interaction) error {
	data, err := info.JSON()
	if err != nil {
		app.Logger.Errorf("格式化回收数据异常 %+v %+v %+v\n", info, rk, err)
		return err
	}
	return pdr.PublishDataLogErr(data, []string{string(rk)}, "lvmio")
}
