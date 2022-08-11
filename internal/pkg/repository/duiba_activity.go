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

func (repo DuiBaActivityRepository) Delete(transaction *repotypes.DeleteDuiBaActivityDO) error {
	return repo.ctx.DB.Model(entity.DuiBaActivity{}).Where("id", transaction.Id).Updates(transaction).Error
}

func (repo DuiBaActivityRepository) GetPageList(by repotypes.GetDuiBaActivityPageDO) (list []entity.DuiBaActivity, total int64, err error) {
	list = make([]entity.DuiBaActivity, 0)
	db := repo.ctx.DB.Table("duiba_activity")
	if by.Type > 0 {
		db.Where("type = ?", by.Type)
	}
	if by.Cid > 0 {
		db.Where("cid = ?", by.Cid)
	}
	if by.Name != "" {
		db.Where("name like ?", "%"+by.Name+"%")
	}
	if by.Statue != 0 {
		db.Where("status = ?", by.Statue)
	}
	if by.ActivityId != "" {
		db.Where("activity_id = ?", by.ActivityId)
	}
	db.Select("*")
	db2 := repo.ctx.DB.Table("(?) as t", db)
	err = db2.Count(&total).
		Offset(by.Offset).
		Limit(by.Limit).
		Order("id desc").
		Find(&list).Error
	if err != nil {
		panic(err)
	}
	return
}

func (repo DuiBaActivityRepository) GetExistOne(do repotypes.GetDuiBaActivityExistDO) (entity.DuiBaActivity, error) {
	ent := entity.DuiBaActivity{}
	db := repo.ctx.DB.Model(ent)
	if do.ActivityId != "" {
		db.Where("activity_id = ?", do.ActivityId)
	}
	if do.NotId != 0 {
		db.Not("id", do.NotId)
	}
	err := db.First(&ent).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return ent, err
	}
	return ent, nil
}
