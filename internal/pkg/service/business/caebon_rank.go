package business

import (
	"errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity/business"
	brepo "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"time"
)

var DefaultCarbonRankService = CarbonRankService{repo: brepo.DefaultCarbonRankRepository}

type CarbonRankService struct {
	repo brepo.CarbonRankRepository
}

// ChangeLikeStatus 排行榜点赞
func (srv CarbonRankService) ChangeLikeStatus(param ChangeLikeStatusParam) (*business.CarbonRankLikeLog, error) {
	start, _ := param.DateType.ParseLastTime()

	rank := srv.repo.FindCarbonRank(brepo.FindCarbonRankBy{
		Pid:        param.Pid,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  start,
	})
	if rank.ID == 0 {
		return nil, errors.New("排行信息不存在")
	}

	like, err := DefaultCarbonRankLikeLogService.ChangeLikeStatus(CarbonRankLikeLogParam{
		Pid:        param.Pid,
		UserId:     param.UserId,
		DateType:   param.DateType,
		ObjectType: param.ObjectType,
		TimePoint:  start,
	})
	if err != nil {
		return nil, err
	}
	likeNumStep := -1
	if like.Status.IsLike() {
		likeNumStep = 1
	}

	rank.LikeNum += likeNumStep

	return like, srv.repo.Save(&rank)
}

// UserRankList 用户碳积分排行榜
func (srv CarbonRankService) UserRankList(param GetUserRankListParam) ([]UserRankInfo, int64, error) {
	start, _ := param.DateType.ParseLastTime()

	list, total, err := srv.repo.GetCarbonRankList(brepo.GetCarbonRankBy{
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
		likeStatusMap[likeStatus.BUserId] = likeStatus.Status.IsLike()
	}

	//需要查询用户信息方法
	userMap := make(map[int64]business.User)

	infoList := make([]UserRankInfo, 0)
	for _, item := range list {
		infoList = append(infoList, UserRankInfo{
			User:    userMap[item.Pid].ShortUser(),
			LikeNum: item.LikeNum,
			IsLike:  likeStatusMap[item.Pid],
			Rank:    item.Rank,
			Value:   item.Value,
		})
	}

	return infoList, total, nil
}

// FindUserRank 获取指定用户的排行榜信息
func (srv CarbonRankService) FindUserRank(param GetMyRankParam) (*UserRankInfo, error) {
	start, _ := param.DateType.ParseLastTime()

	rank := srv.repo.FindCarbonRank(brepo.FindCarbonRankBy{
		Pid:        param.UserId,
		ObjectType: business.RankObjectTypeUser,
		DateType:   param.DateType,
		TimePoint:  start,
	})

	if rank.ID == 0 {
		rank.Rank = 9999
	}

	//需要方法-查询用户信息
	user := business.User{}

	//isLike
	likeStatus, err := DefaultCarbonRankLikeLogService.FindLikeStatus(CarbonRankLikeLogParam{
		Pid:        param.UserId,
		UserId:     param.UserId,
		DateType:   param.DateType,
		ObjectType: business.RankObjectTypeUser,
		TimePoint:  start,
	})
	if err != nil {
		return nil, err
	}

	return &UserRankInfo{
		User:    user.ShortUser(),
		Value:   rank.Value,
		LikeNum: rank.LikeNum,
		Rank:    rank.Rank,
		IsLike:  likeStatus.IsLike(),
	}, nil

}

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
			rankInfo := srv.repo.FindCarbonRank(brepo.FindCarbonRankBy{
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
