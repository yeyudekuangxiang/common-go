package message

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity/message"
)

type (
	IMMessageModel interface {
		Insert(data *message.IM) (*message.IM, error)
		Delete(id int64) error
		Update(data *message.IM) error
		SendMessage(data IMSendMsg) error
		GetMessage(params IMGetMsg) ([]*message.IM, int64, error)
	}

	defaultIMMessageModel struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultIMMessageModel) Insert(data *message.IM) (*message.IM, error) {
	//TODO implement me
	panic("implement me")
}

func (d defaultIMMessageModel) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultIMMessageModel) Update(data *message.IM) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultIMMessageModel) SendMessage(data IMSendMsg) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultIMMessageModel) GetMessage(params IMGetMsg) ([]*message.IM, int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewIMMessageModel(ctx *mioContext.MioContext) IMMessageModel {
	return &defaultIMMessageModel{
		ctx: ctx,
	}
}
