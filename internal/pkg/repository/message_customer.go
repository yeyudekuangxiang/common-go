package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	MessageCustomerModel interface {
		FindAll(params FindMessageParams) ([]entity.UserWebMessage, int64, error)
		HaveReadMessage(params SetHaveReadMessageParams) error
	}

	defaultMessageCustomerModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageCustomerModel) HaveReadMessage(params SetHaveReadMessageParams) error {
	query := d.ctx.DB.Model(&entity.MessageCustomer{}).WithContext(d.ctx.Context).Where("rec_id = ?", params.RecId)

	if len(params.MsgIds) > 1 {
		query = query.Where("message_id in (?)", params.MsgIds)
	} else if params.MsgId != 0 {
		query = query.Where("message_id = ?", params.MsgId)
	}

	return query.Update("status", 2).Error
}

func (d defaultMessageCustomerModel) FindAll(params FindMessageParams) ([]entity.UserWebMessage, int64, error) {
	query := d.ctx.DB.Model(&entity.MessageCustomer{}).WithContext(d.ctx.Context).
		Select("message_customer.*,mc.message_content").
		Joins("left join message_content mc on message_customer.message_id = mc.message_id")
	var resp []entity.UserWebMessage
	var total int64
	if params.RecId != 0 {
		query = query.Where("rec_id = ?", params.RecId)
	}

	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}

	if !params.StartTime.IsZero() {
		query = query.Where("created_at > ?", params.StartTime)
	}

	if !params.EndTime.IsZero() {
		query = query.Where("created_at < ?", params.EndTime)
	}

	if params.Limit != 0 && params.Offset != 0 {
		query = query.Limit(params.Limit).Offset(params.Offset)
	}
	err := query.Count(&total).Order("id asc").Find(&resp).Error

	if err == nil {
		return resp, total, nil
	}

	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	return nil, 0, err
}

func NewMessageCustomerModel(ctx *mioContext.MioContext) MessageCustomerModel {
	return &defaultMessageCustomerModel{
		ctx: ctx,
	}
}
