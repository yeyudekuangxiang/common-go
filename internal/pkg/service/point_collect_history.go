package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultPointCollectHistoryService = PointCollectHistoryService{repo: repository.DefaultPointCollectHistoryRepository}

type PointCollectHistoryService struct {
	repo repository.PointCollectHistoryRepository
}

func (srv PointCollectHistoryService) CreateHistory(param CreateHistoryParam) (*entity.PointCollectHistory, error) {
	t := model.NewTime()
	history := entity.PointCollectHistory{
		OpenId: param.OpenId,
		Type:   param.TransactionType,
		Info:   param.Info,
		Date:   t.Date(),
		Time:   t,
	}
	return &history, srv.repo.Create(&history)
}
