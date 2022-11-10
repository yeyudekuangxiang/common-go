package service

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	repo "mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"time"
)

var DefaultAnswerService = IndexIconService{ctx: context.NewMioContext()}

func NewIndexIconService(ctx *context.MioContext) *IndexIconService {
	return &IndexIconService{
		ctx:  ctx,
		repo: repo.NewIndexIconRepository(ctx),
	}
}

type IndexIconService struct {
	ctx  *context.MioContext
	repo *repo.IndexIconRepository
}

func (srv IndexIconService) DeleteById(dto repotypes.DeleteIndexIconDO) error {
	do := repotypes.DeleteIndexIconDO{
		Id:     dto.Id,
		Status: 2,
	}
	return srv.repo.Delete(&do)
}

func (srv IndexIconService) Create(dto entity.IndexIcon) error {
	err := srv.repo.Create(&entity.IndexIcon{
		Title:     dto.Title,
		RowNum:    dto.RowNum,
		Sort:      dto.Sort,
		Status:    dto.Status,
		IsOpen:    dto.IsOpen,
		Pic:       dto.Pic,
		CreatedAt: time.Now(),
	})
	return err
}

func (srv IndexIconService) Update(dto entity.IndexIcon) error {
	return srv.repo.Update(srv_types.UpdateIndexIconDTO{
		Title:  dto.Title,
		RowNum: dto.RowNum,
		Sort:   dto.Sort,
		Status: dto.Status,
		IsOpen: dto.IsOpen,
		Pic:    dto.Pic,
	})
}

func (srv IndexIconService) Page(dto repotypes.GetIndexIconPageDO) ([]entity.IndexIcon, int64, error) {
	return srv.repo.GetPage(repotypes.GetIndexIconPageDO{
		Offset: dto.Offset,
		Limit:  dto.Limit,
		Title:  dto.Title,
		Status: entity.IndexIconStatusOk,
		IsOpen: dto.IsOpen,
	})
}
