package repository

import (
	"database/sql"
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"time"
)

type (
	CollectionModel interface {
		FindOne(id int64) (*entity.Collection, error)
		Insert(data *entity.Collection) (*entity.Collection, error)
		Delete(id int64) error
		Update(data *entity.Collection) error
		FindAllByOpenId(objType int, openId string, limit, offset int) ([]*entity.Collection, int64, error)
		FindAllByTime(startTime, endTime time.Time, limit, offset int) ([]*entity.Collection, int64, error)
		FindOneByObj(objId int64, objType int, openId string) (*entity.Collection, error)
		Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
	}

	defaultCollectionModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultCollectionModel) Trans(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return d.ctx.DB.Transaction(fc, opts...)
}

func (d defaultCollectionModel) FindOneByObj(objId int64, objType int, openId string) (*entity.Collection, error) {
	var resp entity.Collection
	err := d.ctx.DB.Model(&entity.Collection{}).
		Where("obj_id = ?", objId).
		Where("obj_type = ?", objType).
		Where("open_id= ?", openId).
		First(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (d defaultCollectionModel) FindOne(id int64) (*entity.Collection, error) {
	var resp entity.Collection
	err := d.ctx.DB.Model(&entity.Collection{}).
		First(&resp, id).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (d defaultCollectionModel) Insert(data *entity.Collection) (*entity.Collection, error) {
	err := d.ctx.DB.WithContext(d.ctx.Context).Create(data).Error
	switch err {
	case nil:
		return data, nil
	default:
		return nil, err
	}
}

func (d defaultCollectionModel) Delete(id int64) error {
	result, err := d.FindOne(id)
	if err != nil {
		return err
	}
	return d.ctx.DB.WithContext(d.ctx.Context).Delete(result).Error
}

func (d defaultCollectionModel) Update(data *entity.Collection) error {
	if data.Id == 0 {
		return gorm.ErrPrimaryKeyRequired
	}
	return d.ctx.DB.Save(data).Error
}

func (d defaultCollectionModel) FindAllByOpenId(objType int, openId string, limit, offset int) ([]*entity.Collection, int64, error) {
	var result []*entity.Collection
	var total int64

	query := d.ctx.DB.WithContext(d.ctx.Context).Model(&entity.Collection{}).
		Where("open_id = ?", openId).
		Where("status = ?", 1).
		Where("obj_type = ?", objType)

	if limit != 0 {
		query = query.Limit(limit)
	}

	if offset != 0 {
		query = query.Offset(offset)
	}

	err := query.Count(&total).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (d defaultCollectionModel) FindAllByTime(startTime, endTime time.Time, limit, offset int) ([]*entity.Collection, int64, error) {
	var result []*entity.Collection
	var total int64
	query := d.ctx.DB.WithContext(d.ctx.Context)
	if !startTime.IsZero() {
		query = query.Where("created_at > ?", startTime)
	}

	if !endTime.IsZero() {
		query = query.Where("created_at < ?", endTime)
	}

	if err := query.Where("status = ?", 1).
		Count(&total).
		Limit(limit).
		Offset(offset).
		Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func NewCollectionRepository(ctx *mioContext.MioContext) CollectionModel {
	return &defaultCollectionModel{
		ctx: ctx,
	}
}
