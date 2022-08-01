package business

import (
	"mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultPointLogService = PointLogService{repo: rbusiness.DefaultPointLogRepository}

type PointLogService struct {
	repo rbusiness.PointLogRepository
}

func (srv PointLogService) GetListBy(param GetPointLogListParam) []business.PointLog {
	return srv.repo.GetListBy(rbusiness.GetPointLogListBy{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		OrderBy:   param.OrderBy,
		Type:      param.Type,
	})
}

func (srv PointLogService) CreatePointLog(param CreatePointLogParam) (*business.PointLog, error) {
	log := business.PointLog{
		TransactionId: param.TransactionId,
		BUserId:       param.UserId,
		Type:          param.Type,
		Value:         param.Value,
		OrderId:       param.OrderId,
		Info:          param.Info,
	}
	return &log, srv.repo.Create(&log)
}

func (srv PointLogService) GetUserTotalPointsByUserId(userId int64) rbusiness.GetUserTotalCarbonCredits {
	return srv.repo.GetUserTotalPoints(rbusiness.GetCarbonCreditsLogSortedListBy{UserId: userId})
}
