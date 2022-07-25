package business

import (
	"mio/internal/pkg/model/entity"
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

func (srv PointLogService) GetPointLogInfoList(param GetPointLogInfoListParam) []PointLogInfo {
	ptList := srv.GetListBy(GetPointLogListParam{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		OrderBy:   entity.OrderByList{business.OrderByPointLogCTDESC},
	})

	infoList := make([]PointLogInfo, 0)
	for _, pt := range ptList {
		infoList = append(infoList, PointLogInfo{
			ID:       pt.ID,
			Type:     pt.Type,
			TypeText: pt.Type.Text(),
			TimeStr:  pt.CreatedAt.Format("01.02 15:04:05"),
			Value:    pt.Value,
		})
	}
	return infoList
}
func (srv PointLogService) CreatePointLog(param CreatePointLogParam) (*business.PointLog, error) {
	log := business.PointLog{
		TransactionId: param.TransactionId,
		BUserId:       param.UserId,
		Type:          param.Type,
		Value:         param.Value,
		OrderId:       param.OrderId,
		Info:          business.PointTypeInfo(param.Info),
	}
	return &log, srv.repo.Create(&log)
}

func (srv PointLogService) GetUserTotalPointsByUserId(userId int64) rbusiness.GetUserTotalCarbonCredits {
	return srv.repo.GetUserTotalPoints(rbusiness.GetCarbonCreditsLogSortedListBy{UserId: userId})
}
