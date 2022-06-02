package business

import (
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCarbonRankLikeLogService = CarbonRankLikeLogService{repo: rbusiness.DefaultCarbonRankLikeLogRepository}

type CarbonRankLikeLogService struct {
	repo rbusiness.CarbonRankLikeLogRepository
}

func (srv CarbonRankLikeLogService) FindOrCreateLikeLog(param CarbonRankLikeLogParam) (*ebusiness.CarbonRankLikeLog, error) {

	log := srv.repo.FindLikeLog(rbusiness.FindCarbonRankLikeLogBy{
		Pid:        param.Pid,
		UserId:     param.UserId,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})
	if log.ID != 0 {
		return &log, nil
	}
	log = ebusiness.CarbonRankLikeLog{
		Pid:        param.Pid,
		BUserId:    param.UserId,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  model.Time{Time: param.TimePoint},
		Status:     2,
	}
	return &log, srv.repo.Create(&log)
}
func (srv CarbonRankLikeLogService) ChangeLikeStatus(param CarbonRankLikeLogParam) (*ebusiness.CarbonRankLikeLog, error) {
	log, err := srv.FindOrCreateLikeLog(param)
	if err != nil {
		return nil, err
	}
	log.Status = (log.Status % 2) + 1
	return log, srv.repo.Save(log)
}
func (srv CarbonRankLikeLogService) FindLikeStatus(param CarbonRankLikeLogParam) (int8, error) {
	log := srv.repo.FindLikeLog(rbusiness.FindCarbonRankLikeLogBy{
		Pid:        param.Pid,
		UserId:     param.UserId,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})
	if log.ID != 0 {
		return log.Status, nil
	}
	return 2, nil
}
func (srv CarbonRankLikeLogService) GetLikeLogList(param GetCarbonRankLikeLogListParam) ([]ebusiness.CarbonRankLikeLog, error) {
	list := srv.repo.GetLikeLogList(rbusiness.GetCarbonRankLikeLogListBy{
		PIds:       param.PIds,
		UserId:     param.UserId,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})
	return list, nil
}
