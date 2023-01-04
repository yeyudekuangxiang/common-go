package service

import (
	"mio/config"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/pkg/baidu"
	"mio/pkg/errno"
)

type ReviewService struct {
	ctx          *mioctx.MioContext
	reviewClient *baidu.ReviewClient
}

func NewDefaultReviewClient() *baidu.ReviewClient {
	return &baidu.ReviewClient{
		AccessToken: baidu.NewAccessToken(baidu.AccessTokenConfig{
			RedisClient: app.Redis,
			Prefix:      config.RedisKey.BaiDu,
			AppKey:      config.Config.BaiDuReview.AppKey,
			AppSecret:   config.Config.BaiDuReview.AppSecret,
		}),
	}
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

func (srv ReviewService) ImageReview(param baidu.ImageReviewParam) error {
	err := srv.reviewClient.ImageReview(param)
	if err != nil {
		return err
	}
	return nil
}

func (srv ReviewService) CheckRisk(risk int) error {
	if risk > 2 {
		return errno.ErrCommon.WithMessage("风险等级检测异常，请您稍后再试")
	}
	return nil
}
