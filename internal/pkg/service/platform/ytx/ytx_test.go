package ytx

import (
	"flag"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"testing"
)

func TestGoFunc(t *testing.T) {
	flagConf := flag.String("c", "/Users/yunfeng/Documents/workspace/mp2c-go/config.yaml", "-c")
	initialize.Initialize(*flagConf)

	user := entity.User{
		OpenId: "oy_BA5EsE0mPQvll8eAqPCkBvI8Q",
	}
	bdscene := service.DefaultBdSceneService.FindByCh("yitongxing")
	var options []Options
	options = append(options, WithPoolCode("RP202110251300002"))
	options = append(options, WithSecret("a123456"))
	options = append(options, WithAppId(bdscene.AppId))
	options = append(options, WithDomain(bdscene.Domain))
	ytxService := NewYtxService(mioctx.NewMioContext(), options...)

	_, err := ytxService.SendCoupon(1001, 5.00, user)
	if err != nil {
		app.Logger.Errorf("亿通行发红包失败:%s", err.Error())
		return
	}
}
