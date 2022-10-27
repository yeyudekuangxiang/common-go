package repository

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"time"
)

type (
	MessageModel interface {
		FindOne(id int64) (*entity.Message, error)
		Insert(data *entity.Message) (*entity.Message, error)
		Delete(id int64) error
		Update(data *entity.Message) error
		SendMessage(data SendMessage) error
	}

	defaultMessageModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageModel) SendMessage(data SendMessage) error {
	err := d.ctx.DB.Transaction(func(tx *gorm.DB) error {
		message := entity.Message{
			SendId:    data.SendId,
			RecId:     data.RecId,
			CreatedAt: time.Now(),
		}
		if err := d.ctx.DB.Model(&entity.Message{}).Create(&message).Error; err != nil {
			return err
		}
		messageContent := entity.MessageContent{
			MessageId:      message.Id,
			MessageContent: data.Message,
			CreatedAt:      time.Now(),
		}
		if err := d.ctx.DB.Model(&entity.MessageContent{}).Create(&messageContent).Error; err != nil {
			return err
		}
		messageCustomer := entity.MessageCustomer{
			RecId:     data.RecId,
			MessageId: message.Id,
			CreatedAt: time.Now(),
		}
		if err := d.ctx.DB.Model(&entity.MessageContent{}).Create(&messageCustomer).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (d defaultMessageModel) FindOne(id int64) (*entity.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageModel) Insert(data *entity.Message) (*entity.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageModel) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageModel) Update(data *entity.Message) error {
	//TODO implement me
	panic("implement me")
}

func NewMessageModel(ctx *mioContext.MioContext) MessageModel {
	return &defaultMessageModel{
		ctx: ctx,
	}
}
