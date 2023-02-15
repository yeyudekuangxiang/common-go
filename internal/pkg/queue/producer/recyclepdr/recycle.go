package recyclepdr

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/producer/pdr"
	"mio/internal/pkg/queue/types/message/recyclemsg"
	"mio/internal/pkg/queue/types/routerkey"
)

func Recycle(rk routerkey.RecycleRouterKey, info recyclemsg.IRecycleInfo) error {
	data, err := info.JSON()
	if err != nil {
		app.Logger.Errorf("格式化回收数据异常 %+v %+v %+v\n", info, rk, err)
		return err
	}
	return pdr.PublishDataLogErr(data, []string{string(rk)}, "lvmio")
}
