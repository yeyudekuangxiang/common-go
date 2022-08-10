package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewDuiBaActivityRepository(ctx *context.MioContext) *DuiBaActivityRepository {
	return &DuiBaActivityRepository{
		ctx: ctx,
	}
}

type DuiBaActivityRepository struct {
	ctx *context.MioContext
}

func (repo DuiBaActivityRepository) Save(transaction *entity.DuiBaActivity) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo DuiBaActivityRepository) Create(transaction *entity.DuiBaActivity) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo DuiBaActivityRepository) GetUserChannelPageList(by GetUserChannelPageListBy) (list []entity.DuiBaActivity, total int64) {
	list = make([]entity.DuiBaActivity, 0)

	db := repo.ctx.DB.Table("user_channel")
	if by.Pid > 0 {
		db.Where("pid = ?", by.Pid)
	}
	if by.Cid > 0 {
		db.Where("cid = ?", by.Cid)
	}
	if by.Name != "" {
		db.Where("name like ?", "%"+by.Name+"%")
	}
	db.Select("*")
	db2 := repo.ctx.DB.Table("(?) as t", db)
	err := db2.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("id desc").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}

func (repo DuiBaActivityRepository) GetExistOne(do repotypes.GetDuiBaActivityExistDO) (entity.Banner, error) {
	banner := entity.Banner{}
	db := repo.ctx.DB.Model(banner)
	if do.Name != "" {
		db.Where("name = ?", do.Name)
	}
	if do.ImageUrl != "" {
		db.Or("image_url = ?", do.ImageUrl)
	}
	if do.NotId != 0 {
		db.Not("id", do.NotId)
	}
	err := db.First(&banner).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return banner, err
	}
	return banner, nil
}
