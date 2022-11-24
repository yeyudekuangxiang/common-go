package message

import (
	mioContext "mio/internal/pkg/core/context"
)

type (
	IMMessage interface {
	}

	defaultIMMessage struct {
		ctx     *mioContext.MioContext
		options *imOption
	}

	imOption struct {
	}

	IMOptions func(option *imOption)
)

func NewIMMessageService(ctx *mioContext.MioContext, options ...IMOptions) IMMessage {
	option := &imOption{}
	for i := range options {
		options[i](option)
	}

	return &defaultIMMessage{
		ctx:     ctx,
		options: option,
	}
}

func (d defaultIMMessage) Bind() {

}

func (d defaultIMMessage) Send() {

}

func (d defaultIMMessage) GetByFriend() {

}
