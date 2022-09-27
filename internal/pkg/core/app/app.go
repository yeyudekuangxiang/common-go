package app

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
	"github.com/medivhzhan/weapp/v3"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mio/pkg/wxoa"
)

var (
	// DB 数据库连接
	DB = new(gorm.DB)
	// Logger 日志
	Logger = new(zap.SugaredLogger)
	// Weapp 微信小程序 SDK
	Weapp = new(weapp.Client)
	// WxOa 绿喵服务号 SDK
	WxOa = new(wxoa.WxOA)
	// Redis redis客户端
	Redis = new(redis.Client)
	// OssClient 阿里云oss
	OssClient = new(oss.Client)

	QueueProduct = new(rabbitmq.Publisher)

	RpcService = new(RpcClient)
)

type RpcClient struct {
	CouponRpcSrv couponclient.Coupon
}
