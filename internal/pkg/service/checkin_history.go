package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultCheckinHistoryService = CheckinHistoryService{repo: repository.DefaultCheckinHistoryRepository}

type CheckinHistoryService struct {
	repo repository.CheckinHistoryRepository
}

func (srv CheckinHistoryService) FindCheckHistory(param FindCheckinHistoryParam) (*entity.CheckinHistory, error) {
	history := srv.repo.FindCheckinHistory(repository.FindCheckinHistoryBy{
		OpenId:  param.OpenId,
		OrderBy: param.OrderBy,
	})
	return &history, nil
}
func (srv CheckinHistoryService) FindLastCheckinHistory(openId string) (*entity.CheckinHistory, error) {
	history := srv.repo.FindCheckinHistory(repository.FindCheckinHistoryBy{
		OpenId:  openId,
		OrderBy: entity.OrderByList{entity.OrderByCheckinHistoryTimeDesc},
	})

	return &history, nil
}

func (srv CheckinHistoryService) CreateCheckinHistory(openId string, CheckedNumber int) (*entity.CheckinHistory, error) {
	history := entity.CheckinHistory{
		OpenId:        openId,
		CheckedNumber: CheckedNumber,
		Time:          model.NewTime(),
	}
	return &history, srv.repo.Create(&history)
}
