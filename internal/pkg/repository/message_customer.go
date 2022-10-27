package repository

import mioContext "mio/internal/pkg/core/context"

type (
	MessageCustomerModel interface {
	}

	defaultMessageCustomerModel struct {
		ctx *mioContext.MioContext
	}
)

func NewMessageCustomerModel(ctx *mioContext.MioContext) MessageCustomerModel {
	return &defaultMessageCustomerModel{
		ctx: ctx,
	}
}
