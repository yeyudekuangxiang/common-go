package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
	sdk "github.com/sensorsdata/sa-sdk-go"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxapp"
	"gitlab.miotech.com/miotech-application/backend/common-go/wxoa"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activityclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/carbonpk/carbonpkclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/platform/cmd/rpc/platformclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/point/cmd/rpc/pointclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/tokencenter/cmd/rpc/tokencenterclient"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/user/cmd/rpc/userclient"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"reflect"
)

var (
	// DB 数据库连接
	DB         = new(gorm.DB)
	ActivityDB = new(gorm.DB)
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

	SensorsClient = new(sdk.SensorsAnalytics)
)

type RpcClient struct {
	CarbonPkRpcSrv    carbonpkclient.Carbonpk
	UserRpcSrv        userclient.User
	CouponRpcSrv      couponclient.Coupon
	TokenCenterRpcSrv tokencenterclient.TokenCenter
	PointRpcSrv       pointclient.Point
	ActivityRpcSrv    activityclient.Activity
	PlatformRpcSrv    platformclient.Platform
}

func Ping(ctx context.Context) error {
	//db
	if err := pingDb(ctx, DB); err != nil {
		return err
	}

	//business db
	if err := pingDb(ctx, BusinessDB); err != nil {
		return err
	}

	//redis
	if err := Redis.Ping(ctx).Err(); err != nil {
		return err
	}

	if err := checkNil(ctx, *RpcService); err != nil {
		return err
	}
	_, err := RpcService.PointRpcSrv.Ping(ctx, &pointclient.Request{})
	if err != nil {
		return err
	}
	_, err = RpcService.CouponRpcSrv.Ping(ctx, &couponclient.Request{})
	if err != nil {
		return err
	}
	_, err = RpcService.UserRpcSrv.Ping(ctx, &userclient.Request{})
	if err != nil {
		return err
	}
	_, err = RpcService.ActivityRpcSrv.Ping(ctx, &activityclient.Request{})
	if err != nil {
		return err
	}
	_, err = RpcService.CarbonPkRpcSrv.Ping(ctx, &carbonpkclient.Request{})
	if err != nil {
		return err
	}
	_, err = RpcService.TokenCenterRpcSrv.Ping(ctx, &tokencenterclient.Request{})
	if err != nil {
		return err
	}
	_, err = RpcService.PlatformRpcSrv.Ping(ctx, &platformclient.Request{})
	if err != nil {
		return err
	}
	return nil
}
func pingDb(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("db not init")
	}
	realDb, err := db.DB()
	if err != nil {
		return err
	}
	return realDb.PingContext(ctx)
}
func checkNil(ctx context.Context, client RpcClient) error {
	v := reflect.ValueOf(client)
	vt := reflect.TypeOf(client)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			return fmt.Errorf("%s not init", vt.Field(i).Name)
		}
	}
	return nil
}
