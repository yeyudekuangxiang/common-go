package message

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/message"
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
	query := d.ctx.DB.Model(&message.MessageCustomer{}).WithContext(d.ctx.Context).Where("rec_id = ?", params.RecId)
	if len(params.MsgIds) >= 1 {
		query = query.Where("message_id in (?)", params.MsgIds)
	}
	return query.Update("status", 2).Error
}

func NewMessageCustomerModel(ctx *mioContext.MioContext) MessageCustomerModel {
	return &defaultMessageCustomerModel{
		ctx: ctx,
	}
}
