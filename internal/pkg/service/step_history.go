package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultStepHistoryService = StepHistoryService{repo: repository.DefaultStepHistoryRepository}

type StepHistoryService struct {
	repo repository.StepHistoryRepository
}

func (srv StepHistoryService) FindStepHistory(by FindStepHistoryBy) (*entity.StepHistory, error) {
	step := srv.repo.FindBy(repository.FindStepHistoryBy{
		OpenId:  by.OpenId,
		Day:     by.Day,
		OrderBy: by.OrderBy,
	})
	return &step, nil
}
