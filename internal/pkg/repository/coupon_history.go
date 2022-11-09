package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

var DefaultCouponHistoryRepository = NewCouponHistoryModel(mioContext.NewMioContext())

type (
	CouponHistoryModel interface {
		Insert(data *entity.CouponHistory) (*entity.CouponHistory, error)
		FindOne(id int64) (*entity.CouponHistory, error)
		FindAll(params FindCouponHistoryParams) (*entity.CouponHistory, int64, error)
	}

	defaultCouponHistoryRepository struct {
		ctx *mioContext.MioContext
	}
)

func (m *defaultCouponHistoryRepository) FindAll(params FindCouponHistoryParams) (*entity.CouponHistory, int64, error) {
	resp := &entity.CouponHistory{}
	var total int64
	query := m.ctx.DB.Model(resp).WithContext(m.ctx.Context)
	if params.OpenId != "" {
		query = query.Where("open_id = ?", params.OpenId)
	}

	if params.Types != "" {
		query = query.Where("type = ?", params.Types)
	}

	if !params.StartTime.IsZero() {
		query = query.Where("created_at > ?", params.StartTime)
	}

	if !params.EndTime.IsZero() {
		query = query.Where("created_at < ?", params.EndTime)
	}

	err := query.Count(&total).First(resp).Error

	if err != nil {
		return resp, 0, err
	}

	return resp, total, nil
}

func (m *defaultCouponHistoryRepository) Update(data *entity.CouponHistory) error {
	var result entity.CouponHistory
	err := m.ctx.DB.Model(&result).Where("open_id = ?", data.OpenId).First(&result).Error

	if err != nil {
		return err
	}

	if data.Code != "" {
		result.Code = data.Code
	}

	return m.ctx.DB.Model(&result).Updates(&result).Error
}

func NewCouponHistoryModel(ctx *mioContext.MioContext) CouponHistoryModel {
	return &defaultCouponHistoryRepository{
		ctx: ctx,
	}
}

func (m *defaultCouponHistoryRepository) Insert(data *entity.CouponHistory) (*entity.CouponHistory, error) {
	err := m.ctx.DB.Create(data).Error
	switch err {
	case nil:
		return data, nil
	default:
		return nil, err
	}
}

func (m *defaultCouponHistoryRepository) FindOne(id int64) (*entity.CouponHistory, error) {
	var resp entity.CouponHistory
	err := m.ctx.DB.First(&resp, id).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}
