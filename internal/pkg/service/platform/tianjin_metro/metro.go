package tianjin_metro

import (
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/coupon"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/coupon/cmd/rpc/couponclient"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
)

func NewTianjinMetroService(ctx *context.MioContext) *Service {
	return &Service{
		ctx: ctx,
	}
}

type Service struct {
	ctx *context.MioContext
}

var channelTypes = map[int64]string{
	1073: "天津地铁",
}

func (srv *Service) GetTjMetroTicketStatus(typeId int64, openid string) (*entity.User, error) {
	userInfo, exit, _ := repository.DefaultUserRepository.GetUser(repository.GetUserBy{
		OpenId: openid,
	})
	//判断是否注册绿喵
	if !exit {
		return nil, errno.ErrBindRecordNotFound
	}
	//判断是否指定渠道用户
	_, ok := channelTypes[userInfo.ChannelId]
	if !ok {
		return nil, errno.ErrChannelErr
	}
	//查看是否领取了，没领取满足条件
	couponResp, err := app.RpcService.CouponRpcSrv.FindCoupon(srv.ctx, &couponclient.FindCouponReq{
		UserId:           userInfo.ID,
		CouponCardTypeId: typeId,
	})
	if err != nil {
		return nil, err
	}
	if couponResp.Exist {
		return nil, errno.ErrCouponReceived
	}
	//查看今天是否答题，没答题满足条件
	/*questionsCount, err := quiz.DefaultQuizService.QuestionsCount(openid)
	if err != nil {
		return nil, err
	}
	if questionsCount >= 2 {
		return nil, errno.ErrCommon.WithMessage("不满足答题条件")
	}*/
	return userInfo, nil
}

//测试天津地铁直发

func (srv *Service) SendCoupon(typeId int64, user entity.User) (*coupon.SendCouponV2Resp, error) {
	_, err := srv.GetTjMetroTicketStatus(typeId, user.OpenId)
	if err != nil {
		return nil, err
	}
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
