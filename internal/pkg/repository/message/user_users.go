package message

import (
	"gorm.io/gorm"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/message"
)

type (
	UserUsersModel interface {
		FindOne(channelId string) (*message.UserUsers, error)
		Insert(users message.UserUsers) error
		Delete(channelId string) error
	}

	defaultUserUsersModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultUserUsersModel) Delete(channelId string) error {
	err := d.ctx.DB.WithContext(d.ctx.Context).Where("channel_id = ?", channelId).Delete(&message.UserUsers{}).Error
	switch err {
	case nil:
		return err
	default:
		return err
	}
}

func (d defaultUserUsersModel) Insert(users message.UserUsers) error {
	err := d.ctx.DB.Model(&message.UserUsers{}).WithContext(d.ctx.Context).Create(&users).Error
	switch err {
	case nil:
		return err
	default:
		return err
	}
}

func (d defaultUserUsersModel) FindOne(chanelId string) (*message.UserUsers, error) {
	var resp message.UserUsers
	err := d.ctx.DB.Model(&message.UserUsers{}).WithContext(d.ctx.Context).Where("channel_id = ?", chanelId).First(&resp).Error

	switch err {
	case nil:
		return &resp, err
	case gorm.ErrRecordNotFound:
		return &resp, nil
	default:
		return nil, err
	}
}

func NewUserUsersModel(ctx *mioContext.MioContext) UserUsersModel {
	return &defaultUserUsersModel{
		ctx: ctx,
	}
}
