package initialize

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"log"
	"mio/config"
	"mio/internal/pkg/core/app"
	"strings"
)

func InitRpc() {
	app.RpcService = &app.RpcClient{
		CouponRpcSrv: couponclient.NewCoupon(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: strings.Split(config.Config.CouponRpc.Endpoints, ","),
			Target:    config.Config.CouponRpc.Target,
			NonBlock:  config.Config.CouponRpc.NonBlock,
			Timeout:   config.Config.CouponRpc.Timeout,
		})),
	}
	log.Println("初始化rpc服务成功...")
}
