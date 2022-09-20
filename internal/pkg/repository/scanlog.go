package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type ScanLogRepository struct {
	ctx *context.MioContext
}

func NewScanLogRepository(ctx *context.MioContext) *ScanLogRepository {
	return &ScanLogRepository{ctx: ctx}
}

func (repo ScanLogRepository) Save(log *entity.ScanLog) error {
	return repo.ctx.DB.Save(log).Error
}

func (repo ScanLogRepository) Create(log *entity.ScanLog) error {
	return repo.ctx.DB.Create(log).Error
}

func (repo ScanLogRepository) FindByHash(hash string) (*entity.ScanLog, bool, error) {
	log := entity.ScanLog{}
	err := repo.ctx.DB.Where("hash = ?", hash).Take(&log).Error
	if err == nil {
		return &log, true, nil
	}

	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
}
