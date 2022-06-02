package business

import (
	"mio/internal/pkg/model"
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCarbonRankLikeNumService = CarbonRankLikeNumService{repo: rbusiness.DefaultCarbonRankLikeNumRepository}

type CarbonRankLikeNumService struct {
	repo rbusiness.CarbonRankLikeNumRepository
}

func (srv CarbonRankLikeNumService) GetLikeNumList(param GetCarbonRankLikeNumListParam) ([]ebusiness.CarbonRankLikeNum, error) {
	likeNumList := srv.repo.GetLikeNumList(rbusiness.GetCarbonRankLikeNumListBy{
		PIds:       param.PIds,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})
	return likeNumList, nil
}
func (srv CarbonRankLikeNumService) FindLikeNum(param FindCarbonRankLikeNumParam) (*ebusiness.CarbonRankLikeNum, error) {
	likeNum := srv.repo.FindLikeNum(rbusiness.FindCarbonRankLikeNumBy{
		Pid:        param.Pid,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})
	return &likeNum, nil
}
func (srv CarbonRankLikeNumService) FindOrCreateLikeNum(param CreateCarbonRankLikeNumParam) (*ebusiness.CarbonRankLikeNum, error) {
	likeNum := srv.repo.FindLikeNum(rbusiness.FindCarbonRankLikeNumBy{
		Pid:        param.Pid,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})
	if likeNum.ID != 0 {
		return &likeNum, nil
	}

	likeNum = ebusiness.CarbonRankLikeNum{
		Pid:        param.Pid,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		LikeNum:    0,
		TimePoint:  model.Time{Time: param.TimePoint},
	}
	return &likeNum, srv.repo.Create(&likeNum)
}
func (srv CarbonRankLikeNumService) UpdateLikeNum(param UpdateCarbonRankLikeNumParam) (*ebusiness.CarbonRankLikeNum, error) {
	likeNum, err := srv.FindOrCreateLikeNum(CreateCarbonRankLikeNumParam{
		Pid:        param.Pid,
		ObjectType: param.ObjectType,
		DateType:   param.DateType,
		TimePoint:  param.TimePoint,
	})

	if err != nil {
		return nil, err
	}
	likeNum.LikeNum++
	return likeNum, srv.repo.Save(likeNum)
}
