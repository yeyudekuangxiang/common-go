package business

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	brepo "mio/internal/pkg/repository/business"
)

var DefaultCarbonCreditsLogService = CarbonCreditsLogService{repo: brepo.DefaultCarbonCreditsLogRepository}

type CarbonCreditsLogService struct {
	repo brepo.CarbonCreditsLogRepository
}

func (srv CarbonCreditsLogService) GetCarbonCreditLogList(param GetCarbonCreditLogListParam) []business.CarbonCreditsLog {
	return srv.repo.GetListBy(brepo.GetCarbonCreditsLogListBy{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		OrderBy:   param.OrderBy,
		Type:      param.Type,
	})
}
func (srv CarbonCreditsLogService) GetCarbonCreditLogInfoList(param GetCarbonCreditLogInfoListParam) []CarbonCreditLogInfo {
	cclList := srv.GetCarbonCreditLogList(GetCarbonCreditLogListParam{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		OrderBy:   entity.OrderByList{business.OrderByCarbonCreditsLogCtDesc},
	})

	infoList := make([]CarbonCreditLogInfo, 0)
	for _, ccl := range cclList {
		infoList = append(infoList, CarbonCreditLogInfo{
			ID:       ccl.ID,
			Type:     ccl.Type,
			TypeText: ccl.Type.Text(),
			TimeStr:  ccl.CreatedAt.Format("01.02 15:04:05"),
			Value:    ccl.Value,
		})
	}
	return infoList
}
