package message

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/message"
)

type (
	MessageContentModel interface {
		FindContentByMsgIds(messageId []int64) []message.MessageContent
	}

	defaultMessageContentModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageContentModel) FindContentByMsgIds(messageIds []int64) []message.MessageContent {
	var resp []message.MessageContent
	err := d.ctx.DB.Model(&message.MessageContent{}).WithContext(d.ctx.Context).
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
