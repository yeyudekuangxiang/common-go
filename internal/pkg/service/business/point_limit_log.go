package business

import (
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultPointLimitLogService = PointLimitLogService{repo: rbusiness.DefaultPointLimitLogRepository}

type PointLimitLogService struct {
	repo rbusiness.PointLimitLogRepository
}

func (srv PointLimitLogService) FindLimitLog(param FindPointLimitLogParam) (*ebusiness.PointLimitLog, error) {
	log := srv.repo.FindLimitLog(rbusiness.FindPointLimitLogBy{
		TimePoint: param.TimePoint,
		Type:      param.Type,
		UserId:    param.UserId,
	})
	return &log, nil
}
func (srv PointLimitLogService) UpdateOrCreateLimitLog(param UpdateOrCreatePointLimitLogParam) (*ebusiness.PointLimitLog, error) {
	log := srv.repo.FindLimitLog(rbusiness.FindPointLimitLogBy{
		TimePoint: param.TimePoint,
		Type:      param.Type,
		UserId:    param.UserId,
	})
	if log.ID != 0 {
		log.CurrentCount += 1
		log.CurrentValue += param.AddCurrentValue
		return &log, srv.repo.Save(&log)
	}

	log = ebusiness.PointLimitLog{
		Type:         param.Type,
		TimePoint:    model.Time{Time: param.TimePoint},
		BUserId:      param.UserId,
		CurrentValue: param.AddCurrentValue,
		CurrentCount: 1,
	}
	return &log, srv.repo.Create(&log)
}
