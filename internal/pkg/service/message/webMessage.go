package message

import mioContext "mio/internal/pkg/core/context"

type (
	WebMessage interface {
		SendMessage(sendId, recId int64, tempId int64) error
	}

	defaultWebMessage struct {
		ctx *mioContext.MioContext
	}
)

func (d defaultWebMessage) SendMessage(sendId, recId int64, tempId int64) error {
	//TODO implement me
	panic("implement me")
}

func (d defaultWebMessage) getTemplate(tempId int64) {

}

func NewWebMessageService(ctx *mioContext.MioContext) WebMessage {
	return &defaultWebMessage{
		ctx: ctx,
	}
}
