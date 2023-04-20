package initialize

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activityclient"
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
		//CarbonPkRpcSrv: carbonpkclient.NewCarbonpk(zrpc.MustNewClient(zrpc.RpcClientConf{
		//	Endpoints: endpoints(config.Config.ActivityCarbonPkRpc.Endpoints),
		//	Target:    config.Config.ActivityCarbonPkRpc.Target,
		//	NonBlock:  config.Config.ActivityCarbonPkRpc.NonBlock,
		//	Timeout:   config.Config.ActivityCarbonPkRpc.Timeout,
		//})),
		UserRpcSrv: userclient.NewUser(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: endpoints(config.Config.UserRpc.Endpoints),
			Target:    config.Config.UserRpc.Target,
			NonBlock:  config.Config.UserRpc.NonBlock,
			Timeout:   config.Config.UserRpc.Timeout,
		})),
		//CouponRpcSrv: couponclient.NewCoupon(zrpc.MustNewClient(zrpc.RpcClientConf{
		//	Endpoints: endpoints(config.Config.CouponRpc.Endpoints),
		//	Target:    config.Config.CouponRpc.Target,
		//	NonBlock:  config.Config.CouponRpc.NonBlock,
		//	Timeout:   config.Config.CouponRpc.Timeout,
		//})),
		TokenCenterRpcSrv: tokencenterclient.NewTokenCenter(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: endpoints(config.Config.TokenCenterRpc.Endpoints),
			Target:    config.Config.TokenCenterRpc.Target,
			NonBlock:  config.Config.TokenCenterRpc.NonBlock,
			Timeout:   config.Config.TokenCenterRpc.Timeout,
		})),
		PointRpcSrv: pointclient.NewPoint(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: endpoints(config.Config.PointRpc.Endpoints),
			Target:    config.Config.PointRpc.Target,
			NonBlock:  config.Config.PointRpc.NonBlock,
			Timeout:   config.Config.PointRpc.Timeout,
		})),
		ActivityRpcSrv: activityclient.NewActivity(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: endpoints(config.Config.ActivityRpc.Endpoints),
			Target:    config.Config.ActivityRpc.Target,
			NonBlock:  config.Config.ActivityRpc.NonBlock,
			Timeout:   config.Config.ActivityRpc.Timeout,
		})),
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
