package repository

import mioContext "mio/internal/pkg/core/context"

type (
	MessageContentModel interface {
	}

	defaultMessageContentModel struct {
		ctx *mioContext.MioContext
	}
)

func NewMessageContentModel(ctx *mioContext.MioContext) MessageContentModel {
	return &defaultMessageContentModel{
		ctx: ctx,
	}
}
