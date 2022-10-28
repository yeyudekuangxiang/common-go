package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	MessageContentModel interface {
		FindContentByMsgIds(messageId []int64) []entity.MessageContent
	}

	defaultMessageContentModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageContentModel) FindContentByMsgIds(messageIds []int64) []entity.MessageContent {
	var resp []entity.MessageContent
	err := d.ctx.DB.Model(&entity.MessageContent{}).WithContext(d.ctx.Context).
		Where("message_id in (?)", messageIds).Find(&resp).Error

	if err == nil {
		return resp
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	return nil
}

func NewMessageContentModel(ctx *mioContext.MioContext) MessageContentModel {
	return &defaultMessageContentModel{
		ctx: ctx,
	}
}
