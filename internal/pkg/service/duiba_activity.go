package service

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
)

var DefaultDuiBaActivityService = DuiBaActivityService{repo: repository.DuiBaActivityRepository{}}

type DuiBaActivityService struct {
	repo repository.DuiBaActivityRepository
}

func (srv DuiBaActivityService) FindActivity(activityId string) (*entity.DuiBaActivity, error) {
	activity := entity.DuiBaActivity{}
	err := app.DB.Where("activity_id = ?", activityId).First(&activity).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &activity, nil
}

func (srv DuiBaActivityService) Create(dto srv_types.CreateDuiBaActivityDTO) error {
	//判断名称和图片是否存在
	banner, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		Name: dto.Name})
	if err != nil {
		return err
	}
	if banner.ID != 0 {
		return errors.New("banner名称或图片已存在")
	}
	bannerDo := entity.DuiBaActivity{
		CreatedAt: model.NewTime(),
		UpdatedAt: model.NewTime()}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return err
	}
	return srv.repo.Create(&bannerDo)
}

/*
func (srv DuiBaActivityService) GetBannerPageList(dto srv_types.GetPageBannerDTO) ([]entity.Banner, int64, error) {
	bannerDo := repotypes.GetBannerPageDO{OrderBy: entity.OrderByList{entity.OrderByBannerSortAsc}}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return nil, 0, err
	}
	list, total, err := srv.repo.Page(bannerDo)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
*/
