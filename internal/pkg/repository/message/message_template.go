package message

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/message"
)

type (
	MessageTemplateModel interface {
		FindOne(key string) (*message.MessageTemplate, error)
		Update(data *message.MessageTemplate) error
		Delete(key string) error
		Create(data *message.MessageTemplate) error
	}

	defaultMessageTemplateModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageTemplateModel) FindOne(key string) (*message.MessageTemplate, error) {
	var resp message.MessageTemplate
	err := d.ctx.DB.WithContext(d.ctx.Context).Where("key = ?", key).Where("status = ?", 1).Take(&resp).Error

	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, entity.ErrNotFount
	default:
		return nil, err
	}
}

func (d defaultMessageTemplateModel) Update(data *message.MessageTemplate) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageTemplateModel) Delete(key string) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageTemplateModel) Create(data *message.MessageTemplate) error {
	//TODO implement me
	panic("implement me")
}

func NewMessageTemplateModel(ctx *mioContext.MioContext) MessageTemplateModel {
	return &defaultMessageTemplateModel{
		ctx: ctx,
	}
}
