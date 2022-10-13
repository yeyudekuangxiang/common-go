package service

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
)

var DefaultDuiBaPointAddLogService = DuiBaPointAddLogService{repo: repository.DefaultDuiBaPointAddLogRepository}

type DuiBaPointAddLogService struct {
	repo repository.DuiBaPointAddLogRepository
}

func (srv DuiBaPointAddLogService) FindBy(by FindDuiBaPointAddLogBy) (*entity.DuiBaPointAddLog, error) {
	log := srv.repo.FindBy(repository.FindDuiBaPointAddLogBy{
		OrderNum: by.OrderNum,
	})
	return &log, nil
}
func (srv DuiBaPointAddLogService) CreateLog(add CreateDuiBaPointAddLog) (*entity.DuiBaPointAddLog, error) {
	log := entity.DuiBaPointAddLog{}
	if err := util.MapTo(add, &log); err != nil {
		return nil, err
	}

	return &log, srv.repo.Create(&log)
}
func (srv DuiBaPointAddLogService) UpdateLogTransaction(logId int64, transactionId string) error {
	log := srv.repo.FindByID(logId)
	if log.ID == 0 {
		return errors.New("log 不存在")
	}
	log.TransactionId = transactionId
	return srv.repo.Save(&log)
}
