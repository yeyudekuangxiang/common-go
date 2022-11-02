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
		GetMessage(params FindMessageParams) ([]entity.UserWebMessageV2, int64, error)
		CountAll(params FindMessageParams) (int64, error)
		HaveRead(params FindMessageParams) error
	}

	defaultMessageModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageModel) GetMessage(params FindMessageParams) ([]entity.UserWebMessageV2, int64, error) {
	query := d.ctx.DB.Model(&entity.Message{}).WithContext(d.ctx.Context).
		Select("mcontent.message_id,mcontent.message_content,message.created_at").
		Joins("left join message_content mcontent on message.id = mcontent.message_id").
		Joins("left join message_customer mcustomer on message.id = mcustomer.message_id")
	var resp []entity.UserWebMessageV2
	var total int64

	if params.RecId != 0 {
		query = query.Where("mcustomer.rec_id = ?", params.RecId)
	}

	if params.Status != 0 {
		query = query.Where("mcustomer.status = ?", params.Status)
	}

	if params.Type != 0 {
		query = query.Where("message.type = ?", params.Type)
	} else if len(params.Types) >= 1 {
		query = query.Where("message.type in (?)", params.Types)
	}

	if !params.StartTime.IsZero() {
		query = query.Where("mcustomer.created_at > ?", params.StartTime)
	}

	if !params.EndTime.IsZero() {
		query = query.Where("mcustomer.created_at < ?", params.EndTime)
	}

	if params.Limit != 0 && params.Offset != 0 {
		query = query.Limit(params.Limit).Offset(params.Offset)
	}

	err := query.Count(&total).Order("mcustomer.id asc").Find(&resp).Error

	if err == nil {
		return resp, total, nil
	}

	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	return nil, 0, err
}

func (d defaultMessageModel) CountAll(params FindMessageParams) (int64, error) {
	query := d.ctx.DB.Model(&entity.Message{}).WithContext(d.ctx.Context).
		Joins("left join message_customer mcustomer on message.id = mcustomer.message_id")

	var total int64
	if len(params.MessageIds) > 0 {
		query = query.Where("message.id in (?)", params.MessageIds)
	}

	if params.Type != 0 {
		query = query.Where("message.type = ?", params.Type)
	} else if len(params.Types) >= 1 {
		query = query.Where("message.type in (?)", params.Types)
	}

	if params.RecId != 0 {
		query = query.Where("mcustomer.rec_id = ?", params.RecId)
	}

	if params.Status != 0 {
		query = query.Where("mcustomer.status = ?", params.Status)
	}

	if !params.StartTime.IsZero() {
		query = query.Where("message.created_at > ?", params.StartTime)
	}

	if !params.EndTime.IsZero() {
		query = query.Where("message.created_at < ?", params.EndTime)
	}

	err := query.Count(&total).Error
	if err == nil {
		return total, nil
	}
	return 0, err
}

func (d defaultMessageModel) HaveRead(params FindMessageParams) error {
	query := d.ctx.DB.Model(&entity.MessageCustomer{}).WithContext(d.ctx.Context)

	if len(params.MessageIds) > 0 {
		query = query.Where("message_id in (?)", params.MessageIds)
	}

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

	return query.Update("status", 2).Error
}

func (d defaultMessageModel) SendMessage(params SendMessage) error {
	err := d.ctx.DB.Transaction(func(tx *gorm.DB) error {
		message := entity.Message{
			SendId:    params.SendId,
			RecId:     params.RecId,
			Type:      params.Type,
			CreatedAt: time.Now(),
		}
		if err := d.ctx.DB.Model(&entity.Message{}).Create(&message).Error; err != nil {
			return err
		}
		messageContent := entity.MessageContent{
			MessageId:      message.Id,
			MessageContent: params.Message,
			CreatedAt:      time.Now(),
		}
		if err := d.ctx.DB.Model(&entity.MessageContent{}).Create(&messageContent).Error; err != nil {
			return err
		}
		messageCustomer := entity.MessageCustomer{
			RecId:     params.RecId,
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

func (d defaultMessageModel) Insert(params *entity.Message) (*entity.Message, error) {
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
