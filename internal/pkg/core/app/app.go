package app

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activityclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mio/pkg/wxapp"
	"mio/pkg/wxoa"
)

var (
	// DB 数据库连接
	DB         = new(gorm.DB)
	BusinessDB = new(gorm.DB)
	// Logger 日志
	Logger = new(zap.SugaredLogger)
	// OriginLogger 日志
	OriginLogger = new(zap.Logger)
	// Weapp 微信小程序 SDK
	Weapp = new(wxapp.Client)
	// WxOa 绿喵服务号 SDK
	WxOa = new(wxoa.WxOA)
	// Redis redis客户端
	Redis = new(redis.Client)
	// OssClient 阿里云oss
	OssClient = new(oss.Client)
	// STSClient 阿里云sts
	STSClient = new(sts.Client)

	QueueProduct = new(rabbitmq.Publisher)

	RpcService = new(RpcClient)
)

type RpcClient struct {
	CouponRpcSrv      couponclient.Coupon
	TokenCenterRpcSrv tokencenterclient.TokenCenter
	ActivityRpcSrv    activityclient.Activity
}
