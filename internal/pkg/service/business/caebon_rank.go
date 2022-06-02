package business

import (
	"mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/util/timeutils"
	"time"
)

var DefaultCarbonRankService = CarbonRankService{}

type CarbonRankService struct {
}

func (srv CarbonRankService) ParseDateType(dateType business.RankDateType) (time.Time, time.Time) {
	t := timeutils.Now()
	var start, end time.Time
	switch dateType {
	case business.RankDateTypeDay:
		start = t.AddDay(-1).StartOfDay().Time
		end = t.AddDay(-1).EndOfDay().Time
	case business.RankDateTypeWeek:
		start = t.AddWeek(-1).StartOfWeek().Time
		end = t.AddWeek(-1).EndOfWeek().Time
	case business.RankDateTypeMonth:
		start = t.AddMonth(-1).StartOfMonth().Time
		end = t.AddMonth(-1).EndOfMonth().Time
	}
	return start, end
}

// UserRankList 用户碳积分排行榜
func (srv CarbonRankService) UserRankList(param GetUserRankListParam) ([]UserRankInfo, int64, error) {
	start, end := srv.ParseDateType(param.DateType)

	list, total, err := DefaultCarbonCreditsLogService.GetUserCarbonRank(GetUserCarbonRankParam{
		StartTime: start,
		EndTime:   end,
		CompanyId: param.CompanyId,
		Limit:     param.Limit,
		Offset:    param.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	userIds := make([]int64, 0)
	for _, item := range list {
		userIds = append(userIds, item.UserId)
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
			User:    userMap[item.UserId],
			LikeNum: likeNumMap[item.UserId],
			IsLike:  likeStatusMap[item.UserId],
			Value:   item.Value,
		})
	}

	return infoList, total, nil
}
