package message

import (
	"errors"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/message"
	"mio/pkg/errno"
	"strconv"
	"strings"
)

type (
	WebMessage interface {
		SendMessage(params SendWebMessage) error
		GetMessageCount(params GetWebMessageCount) (GetWebMessageCountResp, error)
		GetMessage(params GetWebMessage) ([]*GetWebMessageResp, int64, error)
		SetHaveRead(params SetHaveReadMessage) error
		GetTemplate(key string) string
		GetTemplateInfo(key string) (*GetMessageTemplate, error)
	}

	defaultWebMessage struct {
		ctx             *mioContext.MioContext
		message         message.MessageModel
		messageCustomer message.MessageCustomerModel
		messageContent  message.MessageContentModel
		template        message.MessageTemplateModel
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

	WMOptions func(option *webMessageOption)
)

func (d defaultWebMessage) SetHaveRead(params SetHaveReadMessage) error {
	err := d.messageCustomer.HaveReadMessage(message.SetHaveReadMessageParams{
		MsgIds: params.MsgIds,
		RecId:  params.RecId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d defaultWebMessage) GetMessageCount(params GetWebMessageCount) (GetWebMessageCountResp, error) {
	res := GetWebMessageCountResp{}

	total, err := d.message.CountAll(message.FindMessageParams{
		RecId:  params.RecId,
		Status: 1,
	})

	if err != nil {
		return res, errno.ErrCommon
	}

	res.Total = total

	exchangeMsgTotal, err := d.message.CountAll(message.FindMessageParams{
		RecId:  params.RecId,
		Status: 1,
		Types:  []string{"1", "2", "3"},
	})

	if err != nil {
		return res, errno.ErrCommon
	}

	res.ExchangeMsgTotal = exchangeMsgTotal

	systemMsgTotal, err := d.message.CountAll(message.FindMessageParams{
		RecId:  params.RecId,
		Status: 1,
		Types:  []string{"4", "5", "6", "7", "8", "9", "10", "11", "12"},
	})
	if err != nil {
		return res, errno.ErrCommon
	}

	res.SystemMsgTotal = systemMsgTotal

	return res, nil
}

func (d defaultWebMessage) GetMessage(params GetWebMessage) ([]*GetWebMessageResp, int64, error) {
	msgList, total, err := d.message.GetMessage(message.FindMessageParams{
		RecId:  params.UserId,
		Status: params.Status,
		Types:  params.Types,
		Limit:  params.Limit,
		Offset: params.Offset,
	})

	if err != nil {
		return nil, 0, err
	}

	l := len(msgList)

	result := make([]*GetWebMessageResp, l)

	if l == 0 {
		return result, 0, nil
	}

	uKeyMap := make(map[int64]struct{}, l+1) // 发送者id map

	for i, item := range msgList {
		one := &GetWebMessageResp{
			Id:             item.Id,
			MessageContent: item.MessageContent,
			MessageNotes:   item.MessageNotes,
			Type:           item.Type,
			Status:         item.Status,
			CreatedAt:      item.CreatedAt,
			TurnType:       item.TurnType,
			TurnId:         strconv.FormatInt(item.TurnId, 10),
			SendId:         item.SendId,
		}
		result[i] = one
		uKeyMap[item.SendId] = struct{}{}
	}

	//User
	var uIds []int64
	for id := range uKeyMap {
		uIds = append(uIds, id)
	}

	uMap := make(map[int64]entity.ShortUser, len(uIds))
	uList := d.user.GetShortUserListBy(repository.GetUserListBy{UserIds: uIds})
	for _, uItem := range uList {
		uMap[uItem.ID] = uItem
	}

	for _, item := range result {
		if item.SendId == 0 {
			item.User.Nickname = "酷喵圈"
			item.User.AvatarUrl = "https://resources.miotech.com/static/mp2c/user/avatar/oy_BA5Jod6_ItzG6bvmPAX2kRYz8/21a36ea8-b252-406e-881c-1ee97334a594.png"
		} else {
			item.User = uMap[item.SendId]
		}
		////文章组合
		//if item.TurnType == 1 {
		//	item.TurnNotes = tMap[item.ShowId]
		//}
		////评论组合
		//if item.TurnType == 2 {
		//	item.TurnNotes = cMap[item.ShowId]
		//}
		////商品组合
		//if item.TurnType == 3 {
		//	item.TurnNotes = oMap[item.ShowId]
		//}
		////订单组合
		//if item.TurnType == 4 {
		//	item.TurnNotes = gMap[item.ShowId]
		//}
	}

	return result, total, nil
}

func (d defaultWebMessage) SendMessage(param SendWebMessage) error {
	content := d.GetTemplate(param.Key)

	if content == "" {
		return errors.New("模板不存在")
	}

	keys := strings.Split(param.Key, "_")
	if len(keys) > 1 {
		if keys[0] == "reply" {
			obj, _ := d.comment.FindOne(param.TurnId)
			content = strings.ReplaceAll(content, "{0}", obj.Message)
		} else {
			if keys[1] == "topic" {
				obj := d.topic.FindById(param.TurnId)
				content = strings.ReplaceAll(content, "{0}", obj.Title)
			}

			if keys[1] == "comment" {
				obj, _ := d.comment.FindOne(param.TurnId)
				content = strings.ReplaceAll(content, "{0}", obj.Message)
			}
		}
	}

	//入库
	err := d.message.SendMessage(message.SendMessage{
		SendId:       param.SendId,
		RecId:        param.RecId,
		Type:         param.Type,
		TurnType:     param.TurnType,
		TurnId:       param.TurnId,
		MessageNotes: param.MessageNotes,
		Message:      content,
	})
	if err != nil {
		return err
	}

	return nil
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

func (d defaultWebMessage) GetTemplate(key string) string {
	one, err := d.template.FindOne(key)
	if err != nil {
		return ""
	}
	return one.TempContent
}

func (d defaultWebMessage) GetTemplateInfo(key string) (*GetMessageTemplate, error) {
	var one, err = d.template.FindOne(key)
	if err != nil {
		return nil, err
	}
	return &GetMessageTemplate{
		Id:          one.Id,
		Key:         one.Key,
		Type:        one.Type,
		TempContent: one.TempContent,
		CreatedAt:   one.CreatedAt,
		UpdatedAt:   one.UpdatedAt,
	}, nil
}

func WithSendObjId(sendObjId int64) WMOptions {
	return func(option *webMessageOption) {
		option.SendObjID = sendObjId
	}
}

func WithRecObjId(recObjId int64) WMOptions {
	return func(option *webMessageOption) {
		option.RecObjId = recObjId
	}
}

func WithVal(val string) WMOptions {
	return func(option *webMessageOption) {
		option.Val = val
	}
}

func NewWebMessageService(ctx *mioContext.MioContext, options ...WMOptions) WebMessage {
	option := &webMessageOption{}
	for i := range options {
		options[i](option)
	}

	return &defaultWebMessage{
		ctx:             ctx,
		message:         message.NewMessageModel(ctx),
		messageCustomer: message.NewMessageCustomerModel(ctx),
		messageContent:  message.NewMessageContentModel(ctx),
		template:        message.NewMessageTemplateModel(ctx),
		user:            repository.NewUserRepository(),
		topic:           repository.NewTopicModel(ctx),
		comment:         repository.NewCommentModel(ctx),
		options:         option,
	}
}