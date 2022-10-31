package message

import (
	"errors"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
	"strings"
)

type (
	WebMessage interface {
		SendMessage(param sendWebMessage) error
		GetMessageCount(userID int64, status, forType int) (int64, error)
		GetMessage(userID int64, status, forType, limit, offset int) ([]entity.UserWebMessage, int64, error)
	}

	defaultWebMessage struct {
		ctx             *mioContext.MioContext
		message         repository.MessageModel
		messageCustomer repository.MessageCustomerModel
		messageContent  repository.MessageContentModel
		template        repository.MessageTemplateModel
		user            repository.UserRepository
		topic           repository.TopicModel
		comment         repository.CommentModel
		options         *webMessageOption
	}

	webMessageOption struct {
		SendObjID int64  `json:"sendObjId"` // 发送者 object id
		RecObjId  int64  `json:"recObjId"`  // 接受者 object id
		Val       string `json:"val"`
	}

	Options func(option *webMessageOption)
)

func (d defaultWebMessage) GetMessageCount(userID int64, status, forType int) (int64, error) {
	total, err := d.message.CountAll(repository.FindMessageParams{
		RecId:  userID,
		Status: status,
		Type:   forType,
	})

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (d defaultWebMessage) GetMessage(userID int64, status, forType, limit, offset int) ([]entity.UserWebMessage, int64, error) {
	msgList, total, err := d.message.GetMessage(repository.FindMessageParams{
		RecId:  userID,
		Status: status,
		Type:   forType,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	var msgIds []int64
	for _, item := range msgList {
		msgIds = append(msgIds, item.MessageId)
	}

	err = d.message.HaveRead(repository.FindMessageParams{MessageIds: msgIds})
	if err != nil {
		app.Logger.Errorf("Message HaveRead Error:%s", err.Error())
	}

	return msgList, total, nil
}

func (d defaultWebMessage) SendMessage(param sendWebMessage) error {
	sendUser, b, err := d.user.GetUserByID(param.SendId)
	if err != nil {
		return err
	}

	if !b {
		return errno.ErrUserNotFound.WithMessage("发送消息用户不存在")
	}

	_, b, err = d.user.GetUserByID(param.RecId)
	if err != nil {
		return err
	}

	if !b {
		return errno.ErrUserNotFound.WithMessage("发送消息用户不存在")
	}

	content := d.getTemplate(param.Key)

	if content == "" {
		return errors.New("模板不存在")
	}

	content = strings.ReplaceAll(content, "userName", sendUser.Nickname)
	keys := strings.Split(param.Key, "_")
	if len(keys) >= 2 {
		if keys[1] == "topic" {
			content = d.replaceTempForTopic(content, param.RecObjId)
		}

		if keys[1] == "comment" {
			content = d.replaceTempForComment(content, param.RecId)
		}

	} else {
		content = strings.ReplaceAll(content, param.Key, d.options.Val)
	}

	//入库
	err = d.message.SendMessage(repository.SendMessage{
		SendId:  param.SendId,
		RecId:   param.RecId,
		Type:    param.Type,
		Message: content,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d defaultWebMessage) replaceTempForTopic(content string, recObjID int64) string {
	var topic entity.Topic
	var title string
	if d.options.SendObjID != 0 {
		topic = d.topic.FindById(d.options.SendObjID)
		title = topic.Title
		topicRune := []rune(topic.Title)
		if len(topicRune) > 5 {
			title = string(topicRune[0:7]) + "..."
		}
		return strings.ReplaceAll(content, "reTopicTitle", title)
	}

	topic = d.topic.FindById(recObjID)
	title = topic.Title
	topicRune := []rune(topic.Title)
	if len(topicRune) > 5 {
		title = string(topicRune[0:7]) + "..."
	}
	return strings.ReplaceAll(content, "topicTitle", title)
}

func (d defaultWebMessage) replaceTempForComment(content string, recObjID int64) string {
	var comment, recComment *entity.CommentIndex
	var message, recMessage string
	if d.options.SendObjID != 0 {
		comment, _ = d.comment.FindOne(d.options.SendObjID)
		message = comment.Message
		messageRune := []rune(message)
		if len(messageRune) > 5 {
			message = string(messageRune[0:5]) + "..."
		}
		content = strings.ReplaceAll(content, "reComment", message)
	}

	recComment, _ = d.comment.FindOne(recObjID)
	recMessage = recComment.Message
	messageRune := []rune(recMessage)
	if len(messageRune) > 5 {
		recMessage = string(messageRune[0:5]) + "..."
	}

	content = strings.ReplaceAll(content, "comment", recMessage)

	return content
}

func (d defaultWebMessage) getTemplate(key string) string {
	one, err := d.template.FindOne(key)
	if err != nil {
		return ""
	}
	return one.TempContent
}

func WithSendObjId(sendObjId int64) Options {
	return func(option *webMessageOption) {
		option.SendObjID = sendObjId
	}
}

func WithRecObjId(recObjId int64) Options {
	return func(option *webMessageOption) {
		option.RecObjId = recObjId
	}
}

func WithVal(val string) Options {
	return func(option *webMessageOption) {
		option.Val = val
	}
}

func NewWebMessageService(ctx *mioContext.MioContext, options ...Options) WebMessage {
	option := &webMessageOption{}
	for i := range options {
		options[i](option)
	}

	return &defaultWebMessage{
		ctx:             ctx,
		message:         repository.NewMessageModel(ctx),
		messageCustomer: repository.NewMessageCustomerModel(ctx),
		messageContent:  repository.NewMessageContentModel(ctx),
		template:        repository.NewMessageTemplateModel(ctx),
		user:            repository.NewUserRepository(),
		topic:           repository.NewTopicRepository(ctx),
		comment:         repository.NewCommentRepository(ctx),
		options:         option,
	}
}
