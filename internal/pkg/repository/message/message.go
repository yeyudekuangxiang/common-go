package message

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/message"
	"time"
)

type (
	MessageModel interface {
		FindOne(id int64) (*message.Message, error)
		Insert(data *message.Message) (*message.Message, error)
		Delete(id int64) error
		Update(data *message.Message) error
		SendMessage(data SendMessage) error
		GetMessage(params FindMessageParams) ([]*message.UserWebMessage, int64, error)
		CountAll(params FindMessageParams) (int64, error)
		HaveRead(params FindMessageParams) error
	}

	defaultMessageModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultMessageModel) GetMessage(params FindMessageParams) ([]*message.UserWebMessage, int64, error) {
	query := d.ctx.DB.Model(&message.Message{}).WithContext(d.ctx.Context).
		Select("message.*,mcontent.message_content,mcontent.message_notes,mcustomer.status").
		Joins("left join message_content mcontent on message.id = mcontent.message_id").
		Joins("left join message_customer mcustomer on message.id = mcustomer.message_id")
	var resp []*message.UserWebMessage
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

	query = query.Count(&total)

	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}

	err := query.Order("mcustomer.status asc,mcustomer.created_at desc").Find(&resp).Error

	if err == nil {
		return resp, total, nil
	}

	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	return nil, 0, err
}

func (d defaultMessageModel) CountAll(params FindMessageParams) (int64, error) {
	query := d.ctx.DB.Model(&message.Message{}).WithContext(d.ctx.Context).
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
	query := d.ctx.DB.Model(&message.MessageCustomer{}).WithContext(d.ctx.Context)

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
		msg := message.Message{
			SendId:    params.SendId,
			RecId:     params.RecId,
			Type:      params.Type,
			TurnType:  params.TurnType,
			TurnId:    params.TurnId,
			CreatedAt: time.Now(),
		}
		if err := d.ctx.DB.Model(&message.Message{}).Create(&msg).Error; err != nil {
			return err
		}
		messageContent := message.MessageContent{
			MessageId:      msg.Id,
			MessageContent: params.Message,
			MessageNotes:   params.MessageNotes,
			CreatedAt:      time.Now(),
		}
		if err := d.ctx.DB.Model(&message.MessageContent{}).Create(&messageContent).Error; err != nil {
			return err
		}
		messageCustomer := message.MessageCustomer{
			RecId:     params.RecId,
			MessageId: msg.Id,
			Status:    1,
			CreatedAt: time.Now(),
		}
		if err := d.ctx.DB.Model(&message.MessageCustomer{}).Create(&messageCustomer).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (d defaultMessageModel) FindOne(id int64) (*message.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageModel) Insert(params *message.Message) (*message.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageModel) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultMessageModel) Update(data *message.Message) error {
	//TODO implement me
	panic("implement me")
}

func NewMessageModel(ctx *mioContext.MioContext) MessageModel {
	return &defaultMessageModel{
		ctx: ctx,
	}
}
