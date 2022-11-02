package repository

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	MessageCustomerModel interface {
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

func NewMessageCustomerModel(ctx *mioContext.MioContext) MessageCustomerModel {
	return &defaultMessageCustomerModel{
		ctx: ctx,
	}
}
