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
func (srv CarbonCreditsLogService) GetUserCarbonRank(param GetUserCarbonRankParam) ([]business.UserCarbonRank, int64, error) {
	return srv.repo.GetUserCarbonRank(brepo.GetUserCarbonRankBy{
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		UserId:    param.UserId,
		CompanyId: param.CompanyId,
		Limit:     param.Limit,
		Offset:    param.Offset,
	})
}
func (srv CarbonCreditsLogService) GetDepartmentCarbonRank(param GetDepartmentCarbonRankParam) ([]business.DepartCarbonRank, int64, error) {
	return srv.repo.GetDepartmentCarbonRank(brepo.GetDepartmentCarbonRankBy{
		StartTime:    param.StartTime,
		EndTime:      param.EndTime,
		DepartmentId: param.DepartmentId,
		CompanyId:    param.CompanyId,
		Limit:        param.Limit,
		Offset:       param.Offset,
	})
}
