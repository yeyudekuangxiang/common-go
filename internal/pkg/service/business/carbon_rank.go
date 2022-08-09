package business

import (
	"errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
	"mio/internal/pkg/util"
	"time"
)

var DefaultCarbonRankService = CarbonRankService{repo: rbusiness.DefaultCarbonRankRepository}

type CarbonRankService struct {
	repo rbusiness.CarbonRankRepository
}

// ChangeLikeStatus 排行榜点赞
func (srv CarbonRankService) ChangeLikeStatus(param ChangeLikeStatusParam) (*business.CarbonRankLikeLog, error) {
	start, _ := param.DateType.ParseLastTime()

	rank := srv.repo.FindCarbonRank(rbusiness.FindCarbonRankBy{
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
		likeStatusMap[likeStatus.Pid] = likeStatus.Status.IsLike()
	}

	userList := DefaultUserService.GetBusinessUserByIds(userIds)
	userMap := make(map[int64]business.User)
	for _, user := range userList {
		userMap[user.ID] = user
	}

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
func (srv CarbonRankService) FindUserRank(param FindUserRankParam) (*UserRankInfo, error) {
	start, _ := param.DateType.ParseLastTime()

	rank := srv.repo.FindCarbonRank(rbusiness.FindCarbonRankBy{
		Pid:        param.UserId,
		ObjectType: business.RankObjectTypeUser,
		DateType:   param.DateType,
		TimePoint:  start,
	})

	if rank.ID == 0 {
		rank.Rank = 9999
	}

	userInfo, err := DefaultUserService.GetBusinessUserById(param.UserId)
	if err != nil {
		return nil, err
	}

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
		User:    userInfo.ShortUser(),
		Value:   rank.Value,
		LikeNum: rank.LikeNum,
		Rank:    rank.Rank,
		IsLike:  likeStatus.IsLike(),
	}, nil

}

// DepartmentRankList 部门碳积分排行榜
func (srv CarbonRankService) DepartmentRankList(param GetDepartmentRankListParam) ([]DepartmentRankInfo, int64, error) {
	start, _ := param.DateType.ParseLastTime()

	list, total, err := srv.repo.GetCarbonRankList(rbusiness.GetCarbonRankBy{
		TimePoint:  start,
		CompanyId:  param.CompanyId,
		Limit:      param.Limit,
		Offset:     param.Offset,
		DateType:   param.DateType,
		ObjectType: business.RankObjectTypeDepartment,
	})

	if err != nil {
		return nil, 0, err
	}

	departmentIds := make([]int64, 0)
	departmentIdsInt := make([]int, 0)
	for _, item := range list {
		departmentIds = append(departmentIds, item.Pid)
		departmentIdsInt = append(departmentIdsInt, int(item.Pid))
	}

	//查询当前用户是否点赞信息
	likeStatusList, err := DefaultCarbonRankLikeLogService.GetLikeLogList(GetCarbonRankLikeLogListParam{
		PIds:       departmentIds,
		UserId:     param.UserId,
		DateType:   param.DateType,
		ObjectType: business.RankObjectTypeDepartment,
		TimePoint:  start,
	})
	if err != nil {
		return nil, 0, err
	}
	likeStatusMap := make(map[int64]bool)
	for _, likeStatus := range likeStatusList {
		likeStatusMap[likeStatus.Pid] = likeStatus.Status.IsLike()
	}

	departmentList := DefaultDepartmentService.GetBusinessDepartmentByIds(departmentIdsInt)
	departmentMap := make(map[int64]business.Department)
	for _, department := range departmentList {
		departmentMap[int64(department.ID)] = department
	}

	infoList := make([]DepartmentRankInfo, 0)
	for _, item := range list {
		infoList = append(infoList, DepartmentRankInfo{
			Department: departmentMap[item.Pid],
			LikeNum:    item.LikeNum,
			IsLike:     likeStatusMap[item.Pid],
			Rank:       item.Rank,
			Value:      item.Value,
		})
	}

	return infoList, total, nil
}

// FindDepartmentRank 获取指定部门的排行榜信息
func (srv CarbonRankService) FindDepartmentRank(param FindDepartmentRankParam) (*DepartmentRankInfo, error) {
	if param.DepartmentId == 0 {
		return &DepartmentRankInfo{
			Department: business.Department{},
			Rank:       9999,
		}, nil
	}

	start, _ := param.DateType.ParseLastTime()

	rank := srv.repo.FindCarbonRank(rbusiness.FindCarbonRankBy{
		Pid:        int64(param.DepartmentId),
		ObjectType: business.RankObjectTypeDepartment,
		DateType:   param.DateType,
		TimePoint:  start,
	})

	if rank.ID == 0 {
		rank.Rank = 9999
	}

	department, err := DefaultDepartmentService.GetBusinessDepartmentById(param.DepartmentId)
	if err != nil {
		return nil, err
	}

	//isLike
	likeStatus, err := DefaultCarbonRankLikeLogService.FindLikeStatus(CarbonRankLikeLogParam{
		Pid:        int64(param.DepartmentId),
		UserId:     param.UserId,
		DateType:   param.DateType,
		ObjectType: business.RankObjectTypeDepartment,
		TimePoint:  start,
	})
	if err != nil {
		return nil, err
	}

	return &DepartmentRankInfo{
		Department: *department,
		Value:      rank.Value,
		LikeNum:    rank.LikeNum,
		Rank:       rank.Rank,
		IsLike:     likeStatus.IsLike(),
	}, nil
}

// InitUserRank 生成用户排名信息
func (srv CarbonRankService) InitUserRank(dateType business.RankDateType) {
	if !util.DefaultLock.Lock(string("InitUserRank_"+dateType), time.Hour*20) {
		app.Logger.Info("20个小时内已经有一个线程初始化过")
		return
	}
	app.Logger.Info("生成用户排行榜", dateType)
	offset := 0
	for {
		list, _, err := DefaultCompanyService.GetCompanyPageList(GetCompanyPageListParam{
			Offset: offset,
			Limit:  100,
		})
		if err != nil {
			app.Logger.Error("查询公司列表异常", err)
			break
		}
		if len(list) == 0 {
			break
		}
		for _, company := range list {
			srv.InitCompanyUserRank(company.ID, dateType)
		}
		offset += 100
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
			return
		}
		if len(list) == 0 {
			return
		}

		for _, item := range list {
			rankInfo := srv.repo.FindCarbonRank(rbusiness.FindCarbonRankBy{
				Pid:        item.UserId,
				ObjectType: business.RankObjectTypeUser,
				DateType:   dateType,
				TimePoint:  start,
			})
			if rankInfo.ID != 0 {
				app.Logger.Error("排行榜数据已存在", companyId, dateType)
				return
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
				return
			}
			rank++
		}

		offset += limit
	}
}

// InitDepartmentRank 生成部门排名信息
func (srv CarbonRankService) InitDepartmentRank(dateType business.RankDateType) {
	if !util.DefaultLock.Lock(string("InitDepartmentRank_"+dateType), time.Hour*20) {
		app.Logger.Info("20个小时内已经有一个线程初始化过")
		return
	}
	app.Logger.Info("生成部门排行榜", dateType)
	offset := 0
	for {
		list, _, err := DefaultCompanyService.GetCompanyPageList(GetCompanyPageListParam{
			Offset: offset,
			Limit:  100,
		})

		if err != nil {
			app.Logger.Error("查询公司列表异常", err)
			break
		}
		if len(list) == 0 {
			break
		}

		for _, company := range list {
			srv.InitCompanyDepartmentRank(company.ID, dateType)
		}
		offset += 100
	}
}

// InitCompanyDepartmentRank 生成公司部门碳积分排行榜
func (srv CarbonRankService) InitCompanyDepartmentRank(companyId int, dateType business.RankDateType) {
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
		list, _, err := DefaultCarbonCreditsLogService.GetActualDepartmentCarbonRank(GetActualDepartmentCarbonRankParam{
			StartTime: start,
			EndTime:   end,
			CompanyId: companyId,
			Limit:     limit,
			Offset:    offset,
		})
		if err != nil {
			app.Logger.Error("生成积分排行榜失败", dateType, companyId, offset, limit, err)
			return
		}
		if len(list) == 0 {
			return
		}

		for _, item := range list {
			//部门等于0跳过
			if item.DepartmentId == 0 {
				continue
			}
			rankInfo := srv.repo.FindCarbonRank(rbusiness.FindCarbonRankBy{
				Pid:        item.DepartmentId,
				ObjectType: business.RankObjectTypeDepartment,
				DateType:   dateType,
				TimePoint:  start,
			})
			if rankInfo.ID != 0 {
				app.Logger.Error("排行榜数据已存在", dateType, companyId)
				return
			}
			rankInfo = business.CarbonRank{
				DateType:   dateType,
				ObjectType: business.RankObjectTypeDepartment,
				Value:      item.Value,
				Rank:       rank,
				Pid:        item.DepartmentId,
				LikeNum:    0,
				TimePoint:  model.Time{Time: start},
			}
			err = srv.repo.Create(&rankInfo)
			if err != nil {
				app.Logger.Error("生成积分排行榜失败", dateType, companyId, offset, limit, err)
				return
			}
			rank++
		}

		offset += limit
	}
}

// DepartmentRankListByCid 根据公司id查询当月排行记录
func (srv CarbonRankService) DepartmentRankListByCid(cid int) ([]business.CarbonRank, error) {
	start, _ := business.RankDateTypeMonth.ParseLastTime()

	list, _, err := srv.repo.GetCarbonRankList(rbusiness.GetCarbonRankBy{
		TimePoint:  start,
		CompanyId:  cid,
		Limit:      1000,
		Offset:     0,
		DateType:   business.RankDateTypeMonth,
		ObjectType: business.RankObjectTypeDepartment,
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}
