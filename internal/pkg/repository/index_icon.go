package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

func NewIndexIconRepository(ctx *context.MioContext) *IndexIconRepository {
	return &IndexIconRepository{
		ctx: ctx,
	}
}

type IndexIconRepository struct {
	ctx *context.MioContext
}

func (repo IndexIconRepository) Save(data *entity.IndexIcon) error {
	return repo.ctx.DB.Save(data).Error
}

func (repo IndexIconRepository) Create(data *entity.IndexIcon) error {
	return repo.ctx.DB.Create(data).Error
}
func (repo IndexIconRepository) CreateInBatches(data []entity.IndexIcon) error {
	return repo.ctx.DB.CreateInBatches(data, len(data)).Error
}

func (repo IndexIconRepository) Delete(by *repotypes.DeleteIndexIconDO) error {
	db := repo.ctx.DB.Model(entity.IndexIcon{})
	if by.Id == 0 {
		return errors.New("id不能为空")
	}
	db.Where("id", by.Id)
	return db.Updates(by).Error
}

func (repo IndexIconRepository) GetOne(do repotypes.GetIndexIconOneDO) (*entity.IndexIcon, bool, error) {
	data := entity.IndexIcon{}
	db := repo.ctx.DB.Model(data)
	if do.ID != 0 {
		db.Where("id = ?", do.ID)
	}
	err := db.First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &data, true, nil
}

func (repo IndexIconRepository) Update(dto srv_types.UpdateIndexIconDTO) error {
	//判断是否存在
	_, exit, err := repo.GetOne(repotypes.GetIndexIconOneDO{
		ID: dto.Id,
	})
	if err != nil {
		return err
	}
	if !exit {
		errno.ErrCommon.WithMessage("金刚位不存在")
	}
	do := entity.IndexIcon{
		UpdatedAt: time.Now()}
	if err := util.MapTo(dto, &do); err != nil {
		return err
	}
	return repo.Save(&do)
}

func (repo IndexIconRepository) GetPage(do repotypes.GetIndexIconPageDO) ([]entity.IndexIcon, int64, error) {
	db := repo.ctx.DB.Model(entity.IndexIcon{})
	if do.Status != 0 {
		db.Where("status = ?", do.Status)
	}
	if do.IsOpen != 0 {
		db.Where("is_open = ?", do.IsOpen)
	}
	if do.Title != "" {
		db.Where("name like ?", "%"+do.Title+"%")
	}
	db.Order("id desc")
	list := make([]entity.IndexIcon, 0)
	var total int64
	return list, total, db.Count(&total).Offset(do.Offset).Limit(do.Limit).Find(&list).Error
}
