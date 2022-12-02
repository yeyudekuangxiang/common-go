package initialize

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activityclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"strings"
)

func InitRpc() {
	app.RpcService = &app.RpcClient{
		CouponRpcSrv: couponclient.NewCoupon(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: endpoints(config.Config.CouponRpc.Endpoints),
			Target:    config.Config.CouponRpc.Target,
			NonBlock:  config.Config.CouponRpc.NonBlock,
			Timeout:   config.Config.CouponRpc.Timeout,
		})),
		TokenCenterRpcSrv: tokencenterclient.NewTokenCenter(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: endpoints(config.Config.TokenCenterRpc.Endpoints),
			Target:    config.Config.TokenCenterRpc.Target,
			NonBlock:  config.Config.TokenCenterRpc.NonBlock,
			Timeout:   config.Config.TokenCenterRpc.Timeout,
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
