package business

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCarbonCreditsLogService = CarbonCreditsLogService{repo: rbusiness.DefaultCarbonCreditsLogRepository}

type CarbonCreditsLogService struct {
	repo rbusiness.CarbonCreditsLogRepository
}

func (srv CarbonCreditsLogService) GetCarbonCreditLogList(param GetCarbonCreditLogListParam) []business.CarbonCreditsLog {
	return srv.repo.GetListBy(rbusiness.GetCarbonCreditsLogListBy{
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
func (srv CarbonCreditsLogService) GetActualUserCarbonRank(param GetActualUserCarbonRankParam) ([]business.UserCarbonRank, int64, error) {
	return srv.repo.GetActualUserCarbonRank(rbusiness.GetActualUserCarbonRankBy{
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		CompanyId: param.CompanyId,
		Limit:     param.Limit,
		Offset:    param.Offset,
	})
}

func (srv CarbonCreditsLogService) GetActualDepartmentCarbonRank(param GetActualDepartmentCarbonRankParam) ([]business.DepartCarbonRank, int64, error) {
	return srv.repo.GetActualDepartmentCarbonRank(rbusiness.GetActualDepartmentCarbonRankBy{
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		CompanyId: param.CompanyId,
		Limit:     param.Limit,
		Offset:    param.Offset,
	})
}
func (srv CarbonCreditsLogService) CreateCarbonCreditLog(param CreateCarbonCreditLogParam) (*business.CarbonCreditsLog, error) {
	log := business.CarbonCreditsLog{
		TransactionId: param.TransactionId,
		BUserId:       param.UserId,
		Type:          param.Type,
		Value:         param.Value,
		Info:          business.CarbonTypeInfo(param.Info),
	}
	return &log, srv.repo.Create(&log)
}
