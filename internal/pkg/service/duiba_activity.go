package service

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"time"
)

type DuiBaActivityService struct {
	ctx  *context.MioContext
	repo *repository.DuiBaActivityRepository
}

func NewDuiBaActivityService(ctx *context.MioContext) *DuiBaActivityService {
	return &DuiBaActivityService{
		ctx:  ctx,
		repo: repository.NewDuiBaActivityRepository(ctx),
	}
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
	//判断是否存在
	banner, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		ActivityId: dto.ActivityId})
	if err != nil {
		return err
	}
	if banner.ID != 0 {
		return errors.New("activityId已存在")
	}
	bannerDo := entity.DuiBaActivity{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return err
	}
	return srv.repo.Create(&bannerDo)
}

func (srv DuiBaActivityService) Update(dto srv_types.UpdateDuiBaActivityDTO) error {
	//判断是否存在
	info, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		Id: dto.Id})
	if err != nil {
		return err
	}
	if info.ID == 0 {
		return errors.New("activityId不存在")
	}
	//是否存在
	one, errInfo := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{ActivityId: dto.ActivityId, NotId: dto.Id})
	if errInfo != nil {
		return errInfo
	}
	if one.ID != 0 {
		return errors.New("activityId已存在")
	}
	do := entity.DuiBaActivity{
		UpdatedAt: time.Now()}
	if err := util.MapTo(dto, &do); err != nil {
		return err
	}
	return srv.repo.Save(&do)
}

func (srv DuiBaActivityService) GetPageList(dto srv_types.GetPageDuiBaActivityDTO) ([]entity.DuiBaActivity, int64, error) {
	bannerDo := repotypes.GetDuiBaActivityPageDO{
		Statue: dto.Status,
	}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return nil, 0, err
	}
	list, total, err := srv.repo.GetPageList(bannerDo)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (srv DuiBaActivityService) Delete(dto srv_types.DeleteDuiBaActivityDTO) error {
	//判断是否存在
	info, err := srv.repo.GetExistOne(repotypes.GetDuiBaActivityExistDO{
		Id: dto.Id})
	if err != nil {
		return err
	}
	if info.ID == 0 {
		return errors.New("activityId不存在")
	}
	do := repotypes.DeleteDuiBaActivityDO{
		UpdatedAt: time.Now(),
		Id:        dto.Id,
		Status:    entity.DuiBaActivityStatusNo,
	}
	return srv.repo.Delete(&do)
}
