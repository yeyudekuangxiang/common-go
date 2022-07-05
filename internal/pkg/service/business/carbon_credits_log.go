package business

import (
	"fmt"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"strconv"
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
		Info:          param.Info,
	}
	return &log, srv.repo.Create(&log)
}
func (srv CarbonCreditsLogService) GetCarbonCreditLogSortedList(param GetCarbonCreditLogSortedListParam) []rbusiness.CarbonCreditsLogSortedList {
	return srv.repo.GetSortedListBy(rbusiness.GetCarbonCreditsLogSortedListBy{
		UserId:    param.UserId,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
	})
}

func (srv CarbonCreditsLogService) GetCarbonCreditLogSortedListByCid(param GetCarbonCreditLogSortedListByCidParam) []rbusiness.CarbonCreditsLogSortedList {
	if param.Cid > 0 {
		userList := DefaultUserService.GetBusinessUserListByCid(param.Cid)
		if len(userList) == 0 || userList[0].ID == 0 {
			return []rbusiness.CarbonCreditsLogSortedList{}
		}
		var ids []int64
		for _, v := range userList {
			ids = append(ids, v.ID)
		}
		return srv.repo.GetSortedListBy(rbusiness.GetCarbonCreditsLogSortedListBy{
			UserIds:   ids,
			StartTime: param.StartTime,
			EndTime:   param.EndTime,
		})
	}
	return []rbusiness.CarbonCreditsLogSortedList{}
}

//test
func (srv CarbonCreditsLogService) FormatCarbonCreditLogList(list []rbusiness.CarbonCreditsLogSortedList) []CarbonCreditsLogSortedListResponse {
	var res []CarbonCreditsLogSortedListResponse
	//查询减碳场景
	companyCarbonSceneList := DefaultCompanyCarbonSceneService.GetBusinessCompanyCarbonSceneListBy(rbusiness.GetCompanyCarbonSceneListBy{Status: 1, BCompanyId: 1})
	//最终组合
	for _, v := range companyCarbonSceneList {
		var item CarbonCreditsLogSortedListResponse
		item.Title = v.Title
		item.Type = v.Type
		item.Icon = v.Icon
		for _, l := range list {
			if v.Type == l.Type {
				item.Total = l.Total
			}
		}
		res = append(res, item)
	}
	//排序
	for k, v := range list {
		for k2, l := range res {
			if v.Type == l.Type {
				t := res[k]
				res[k] = res[k2]
				res[k2] = t
				break
			}
		}
	}
	return res
}

func (srv CarbonCreditsLogService) GetCarbonCreditLogListHistoryBy(by rbusiness.GetCarbonCreditsLogSortedListBy) map[string]CarbonCreditLogListHistoryResponse {
	var CarbonCreditLogListHistoryResponseMap map[string]CarbonCreditLogListHistoryResponse
	CarbonCreditLogListHistoryResponseMap = make(map[string]CarbonCreditLogListHistoryResponse)
	listHistory := srv.repo.GetCarbonCreditsLogListHistory(by)

	for _, v := range listHistory {
		v.Title = v.Type.Text()
		if _, ok := CarbonCreditLogListHistoryResponseMap[v.Month]; !ok {
			var detail []rbusiness.CarbonCreditsLogListHistory
			detail = append(detail, v)
			CarbonCreditLogListHistoryResponseMap[v.Month] = CarbonCreditLogListHistoryResponse{Total: v.Total, Month: v.Month, Detail: detail}
		} else {
			newDetail := append(CarbonCreditLogListHistoryResponseMap[v.Month].Detail, v)
			lastTotal, _ := strconv.ParseFloat(CarbonCreditLogListHistoryResponseMap[v.Month].Total, 64)
			thisTotal, _ := strconv.ParseFloat(v.Total, 64)
			newTotal := strconv.Itoa(int(lastTotal + thisTotal))
			fmt.Println(CarbonCreditLogListHistoryResponseMap[v.Month].Total, lastTotal, lastTotal, lastTotal)
			CarbonCreditLogListHistoryResponseMap[v.Month] = CarbonCreditLogListHistoryResponse{Total: newTotal, Month: v.Month, Detail: newDetail}
		}
	}
	return CarbonCreditLogListHistoryResponseMap
}

func (srv CarbonCreditsLogService) GetUserTotalCarbonCreditsByUserId(userId int64) rbusiness.GetUserTotalCarbonCredits {
	return srv.repo.GetUserTotalCarbonCredits(rbusiness.GetCarbonCreditsLogSortedListBy{UserId: userId})
}

func (srv CarbonCreditsLogService) GetUserTotalCarbonCreditsByCid(cid int) rbusiness.GetUserTotalCarbonCredits {
	if cid <= 0 {
		return rbusiness.GetUserTotalCarbonCredits{}
	}
	userList := DefaultUserService.GetBusinessUserListByCid(cid)
	if len(userList) == 0 || userList[0].ID == 0 {
		return rbusiness.GetUserTotalCarbonCredits{}
	}
	var ids []int64
	for _, v := range userList {
		ids = append(ids, v.ID)
	}
	return srv.repo.GetUserTotalCarbonCredits(rbusiness.GetCarbonCreditsLogSortedListBy{UserIds: ids})
}
