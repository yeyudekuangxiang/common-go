package tianjin_metro

import (
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/coupon"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
)

func NewTianjinMetroService(ctx *context.MioContext) *Service {
	return &Service{
		ctx: ctx,
	}
}

type Service struct {
	ctx *context.MioContext
}

//测试天津地铁直发

func (srv *Service) SendCoupon(typeId int64, user entity.User) (*coupon.SendCouponV2Resp, error) {
	//调用微服务，发地铁券
	SendCouponV2Resp, err := app.RpcService.CouponRpcSrv.SendCouponV2(srv.ctx, &couponclient.SendCouponV2Req{
		UserId:              user.ID,
		ThirdUserId:         user.PhoneNumber,
		CouponCardTypeId:    typeId,
		BizId:               util.UUID(),
		DistributionChannel: "天津地铁答题发电子票",
	})
	if err != nil {
		return nil, err
	}
	return SendCouponV2Resp, nil
}
