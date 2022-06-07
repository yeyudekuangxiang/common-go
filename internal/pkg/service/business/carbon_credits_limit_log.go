package business

import (
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCarbonCreditsLimitLogService = CarbonCreditsLimitLogService{repo: rbusiness.DefaultCarbonCreditsLimitLogRepository}

type CarbonCreditsLimitLogService struct {
	repo rbusiness.CarbonCreditsLimitLogRepository
}

func (srv CarbonCreditsLimitLogService) FindLimitLog(param FindCarbonCreditsLimitLogParam) (*ebusiness.CarbonCreditsLimitLog, error) {
	log := srv.repo.FindLimitLog(rbusiness.FindCarbonCreditsLimitLogBy{
		TimePoint: param.TimePoint,
		Type:      param.Type,
		UserId:    param.UserId,
	})
	return &log, nil
}
func (srv CarbonCreditsLimitLogService) UpdateOrCreateLimitLog(param UpdateOrCreateCarbonCreditsLimitLogParam) (*ebusiness.CarbonCreditsLimitLog, error) {
	log := srv.repo.FindLimitLog(rbusiness.FindCarbonCreditsLimitLogBy{
		TimePoint: param.TimePoint,
		Type:      param.Type,
		UserId:    param.UserId,
	})
	if log.ID != 0 {
		log.CurrentCount += 1
		log.CurrentValue = log.CurrentValue.Add(param.AddCurrentValue)
		return &log, srv.repo.Save(&log)
	}

	log = ebusiness.CarbonCreditsLimitLog{
		Type:         param.Type,
		TimePoint:    model.Time{Time: param.TimePoint},
		BUserId:      param.UserId,
		CurrentValue: param.AddCurrentValue,
		CurrentCount: 1,
	}
	return &log, srv.repo.Create(&log)
}
