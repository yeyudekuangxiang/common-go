package service

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

//var DefaultPointCollectHistoryService = PointCollectHistoryService{repo: repository.NewPointCollectHistoryRepository()}

type PointCollectHistoryService struct {
	ctx  *context.MioContext
	repo *repository.PointCollectHistoryRepository
}

func NewPointCollectHistoryService(ctx *context.MioContext) *PointCollectHistoryService {
	return &PointCollectHistoryService{
		ctx:  ctx,
		repo: repository.NewPointCollectHistoryRepository(ctx),
	}
}

func (srv PointCollectHistoryService) CreateHistory(param CreateHistoryParam) (*entity.PointCollectHistory, error) {
	t := model.NewTime()
	history := entity.PointCollectHistory{
		OpenId: param.OpenId,
		Type:   string(param.TransactionType),
		Info:   param.Info,
		Date:   t.Date(),
		Time:   t,
	}
	return &history, srv.repo.Create(&history)
}
