package repository

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
)

type (
	MessageModel interface {
		FindOne(id int64) (*entity.Message, error)
		Insert(data *entity.Message) (*entity.Message, error)
		Delete(id int64) error
		Update(data *entity.Message) error
	}

	defaultMessageModel struct {
		ctx *mioContext.MioContext
	}
)

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

func NewMessageRepository(ctx *mioContext.MioContext) MessageModel {
	return &defaultMessageModel{
		ctx: ctx,
	}
}
