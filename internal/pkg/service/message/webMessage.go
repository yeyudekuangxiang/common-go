package message

import (
	"errors"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
	"strings"
)

type (
	WebMessage interface {
		SendMessage(sendId, recId int64, key string) error
	}

	defaultWebMessage struct {
		ctx      *mioContext.MioContext
		message  repository.MessageModel
		template repository.MessageTemplateModel
		user     repository.UserRepository
	}
)

func (d defaultWebMessage) SendMessage(sendId, recId int64, key string, objId int64) error {
	sendUser, b1, err := d.user.GetUserByID(sendId)
	if err != nil {
		return err
	}

	if !b1 {
		return errno.ErrUserNotFound.WithMessage("发送消息用户不存在")
	}

	content := d.getTemplate(key)

	if content == "" {
		return errors.New("模板不存在")
	}

}

func (d defaultWebMessage) replaceTemplate(content string, sendId, recId int64, key string, objId int64) string {
	keys := strings.Split(strings.ToLower(key), "_")
	if len(keys) >= 2 {
		if keys[0] == "reply" {

		}
		if keys[0] == "fail" {

		}

	}

}

func (d defaultWebMessage) getTemplate(key string) string {
	one, err := d.template.FindOne(key)
	if err != nil {
		return ""
	}
	return one.TempContent
}

func NewWebMessageService(ctx *mioContext.MioContext) WebMessage {
	return &defaultWebMessage{
		ctx:      ctx,
		message:  repository.NewMessageModel(ctx),
		template: repository.NewMessageTemplateModel(ctx),
		user:     repository.NewUserRepository(),
	}
}
