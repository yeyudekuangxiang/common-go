package business

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"time"
)

var DefaultCarbonRankService = CarbonRankService{}

type CarbonRankService struct {
	repo rbusiness.CarbonRankRepository
}

// UserRankList 用户碳积分排行榜
func (srv CarbonRankService) UserRankList(param GetUserRankListParam) ([]UserRankInfo, int64, error) {
	start, _ := param.DateType.ParseLastTime()

	list, total, err := srv.repo.GetCarbonRankList(rbusiness.GetCarbonRankBy{
		TimePoint:  start,
		CompanyId:  param.CompanyId,
		Limit:      param.Limit,
		Offset:     param.Offset,
		DateType:   param.DateType,
		ObjectType: business.RankObjectTypeUser,
	})

	if err != nil {
		return nil, 0, err
	}

	userIds := make([]int64, 0)
	for _, item := range list {
		userIds = append(userIds, item.Pid)
	}

	//查询点赞数量信息
	likeNumList, err := DefaultCarbonRankLikeNumService.GetLikeNumList(GetCarbonRankLikeNumListParam{
		PIds:       userIds,
		ObjectType: business.RankObjectTypeUser,
		DateType:   param.DateType,
		TimePoint:  start,
	})
	if err != nil {
		return nil, 0, err
	}
	likeNumMap := make(map[int64]int)
	for _, likeNum := range likeNumList {
		likeNumMap[likeNum.Pid] = likeNum.LikeNum
	}

	//查询当前用户是否点赞信息
	likeStatusList, err := DefaultCarbonRankLikeLogService.GetLikeLogList(GetCarbonRankLikeLogListParam{
		PIds:       userIds,
		UserId:     param.UserId,
		DateType:   param.DateType,
		ObjectType: business.RankObjectTypeUser,
		TimePoint:  start,
	})
	if err != nil {
		return nil, 0, err
	}
	likeStatusMap := make(map[int64]bool)
	for _, likeStatus := range likeStatusList {
		likeStatusMap[likeStatus.BUserId] = likeStatus.Status == 1
	}

	//需要查询用户信息方法
	userMap := make(map[int64]business.User)

	infoList := make([]UserRankInfo, 0)
	for _, item := range list {
		infoList = append(infoList, UserRankInfo{
			User:    userMap[item.Pid],
			LikeNum: likeNumMap[item.Pid],
			IsLike:  likeStatusMap[item.Pid],
			Value:   item.Value,
		})
	}

	return infoList, total, nil
}

/*func (srv CarbonRankService) GetMyRank(param GetMyRankParam) (*UserRankInfo, error) {
	start, end := srv.ParseDateType(param.DateType)

	list, _, err := DefaultCarbonCreditsLogService.GetUserCarbonRank(GetUserCarbonRankParam{
		StartTime: start,
		EndTime:   end,
		CompanyId: param.CompanyId,
		Limit:     1,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

}*/

// InitUserRank 生成用户排名信息
func (srv CarbonRankService) InitUserRank(dateType business.RankDateType) {
	if !util.DefaultLock.Lock(string("InitUserRank_"+dateType), time.Hour*20) {
		app.Logger.Info("20个小时内已经有一个线程初始化过")
		return
	}

	companyIds := make([]int, 0)
	for _, companyId := range companyIds {
		srv.InitCompanyUserRank(companyId, dateType)
	}
}

// InitCompanyUserRank 生成公司用户碳积分排行榜
func (srv CarbonRankService) InitCompanyUserRank(companyId int, dateType business.RankDateType) {
	limit := 500
	offset := 0
	rank := 1
	start, end := dateType.ParseLastTime()

	defer func() {
		err := recover()
		if err != nil {
			app.Logger.Error("生成积分排行榜失败", dateType, companyId, offset, limit, err)
		}
	}()

	for {
		list, _, err := DefaultCarbonCreditsLogService.GetActualUserCarbonRank(GetActualUserCarbonRankParam{
			StartTime: start,
			EndTime:   end,
			CompanyId: companyId,
			Limit:     limit,
			Offset:    offset,
		})
		if err != nil {
			app.Logger.Error("生成积分排行榜失败", dateType, companyId, offset, limit, err)
			break
		}
		if len(list) == 0 {
			break
		}

		for _, item := range list {
			rankInfo := srv.repo.FindCarbonRank(rbusiness.FindCarbonRankBy{
				Pid:        item.UserId,
				ObjectType: business.RankObjectTypeUser,
				DateType:   dateType,
				TimePoint:  start,
			})
			if rankInfo.ID != 0 {
				rankInfo.Rank = rank
				err = srv.repo.Save(&rankInfo)
				if err != nil {
					app.Logger.Error("生成积分排行榜失败", dateType, companyId, offset, limit, err)
					break
				}
				continue
			}
			rankInfo = business.CarbonRank{
				DateType:   dateType,
				ObjectType: business.RankObjectTypeUser,
				Value:      item.Value,
				Rank:       rank,
				Pid:        item.UserId,
				LikeNum:    0,
				TimePoint:  model.Time{Time: start},
			}
			err = srv.repo.Create(&rankInfo)
			if err != nil {
				app.Logger.Error("生成积分排行榜失败", dateType, companyId, offset, limit, err)
				break
			}
			rank++
		}

		offset += limit
	}
}
