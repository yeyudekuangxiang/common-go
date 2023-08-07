package initialize

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activityclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/carbonpk/carbonpkclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/platform/cmd/rpc/platformclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/point/cmd/rpc/pointclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/user/cmd/rpc/userclient"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"strings"
)

func InitRpc() {
	app.RpcService = &app.RpcClient{
		CarbonPkRpcSrv:    carbonpkclient.NewCarbonpk(zrpc.MustNewClient(config.Config.ActivityCarbonPkRpc)),
		UserRpcSrv:        userclient.NewUser(zrpc.MustNewClient(config.Config.UserRpc)),
		CouponRpcSrv:      couponclient.NewCoupon(zrpc.MustNewClient(config.Config.CouponRpc)),
		TokenCenterRpcSrv: tokencenterclient.NewTokenCenter(zrpc.MustNewClient(config.Config.TokenCenterRpc)),
		PointRpcSrv:       pointclient.NewPoint(zrpc.MustNewClient(config.Config.PointRpc)),
		ActivityRpcSrv:    activityclient.NewActivity(zrpc.MustNewClient(config.Config.ActivityRpc)),
		PlatformRpcSrv:    platformclient.NewPlatform(zrpc.MustNewClient(config.Config.PlatformRpc)),
	}

	log.Println("初始化rpc服务成功...")
}

func endpoints(str string) []string {
	var edps []string
	if str != "" {
		edps = strings.Split(str, ",")
	}
	return edps
}
