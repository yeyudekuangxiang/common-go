package service

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/baidu"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/util/factory"
	"mio/pkg/errno"
)

type ReviewService struct {
	ctx          *mioctx.MioContext
	reviewClient *baidu.ReviewClient
}

func NewDefaultReviewClient() *baidu.ReviewClient {
	return factory.NewBaiDuReviewFromTokenCenterRpc("baidureview", app.RpcService.TokenCenterRpcSrv)
}

func NewReviewService(mioContext *mioctx.MioContext, reviewClient *baidu.ReviewClient) *ReviewService {
	return &ReviewService{
		ctx:          mioContext,
		reviewClient: reviewClient,
	}
}
func DefaultReviewService() *ReviewService {
	return NewReviewService(mioctx.NewMioContext(), NewDefaultReviewClient())
}

func (srv ReviewService) ReviewImage(param baidu.ImageReviewParam) error {
	resp, err := srv.reviewClient.ImageReview(param)
	if err != nil {
		return errno.ErrCheckErr.WithMessage(fmt.Sprintf("系统错误: %s", err.Error()))
	}
	if !resp.IsSuccess() {
		app.Logger.Infof("review err : image_review param is %v, resp is %v", param, resp)
		return errno.ErrCheckErr.WithMessage(fmt.Sprintf("系统错误: %s", resp.ErrorMsg))
	}

	if resp.ConclusionType == 4 {
		return errno.ErrCheckErr.WithMessage("审核失败")
	}

	if resp.ConclusionType != 1 {
		return errno.ErrCheckErr.WithMessage(resp.Data[0].Msg)
	}

	return nil
}
func (srv ReviewService) ReviewImages(list []baidu.ImageReviewParam) error {
	for _, item := range list {
		if err := srv.ReviewImage(item); err != nil {
			return err
		}
	}
	return nil
}
func (srv ReviewService) CheckRisk(risk int) error {
	if risk > 2 {
		return errno.ErrCommon.WithMessage("风险等级检测异常，请您稍后再试")
	}
	return nil
}
