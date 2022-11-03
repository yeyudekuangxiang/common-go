package message

import (
	"errors"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strings"
)

type (
	WebMessage interface {
		SendMessage(params SendWebMessage) error
		GetMessageCount(params GetWebMessageCount) (GetWebMessageCountResp, error)
		GetMessage(params GetWebMessage) ([]*GetWebMessageResp, int64, error)
		SetHaveRead(params SetHaveReadMessage) error
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

func (d defaultWebMessage) SetHaveRead(params SetHaveReadMessage) error {
	err := d.messageCustomer.HaveReadMessage(repository.SetHaveReadMessageParams{
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

	total, err := d.message.CountAll(repository.FindMessageParams{
		RecId:  params.RecId,
		Status: 1,
	})

	if err != nil {
		return res, errno.ErrCommon
	}

	res.Total = total

	exchangeMsgTotal, err := d.message.CountAll(repository.FindMessageParams{
		RecId:  params.RecId,
		Status: 1,
		Types:  []string{"1", "2", "3"},
	})

	if err != nil {
		return res, errno.ErrCommon
	}

	res.ExchangeMsgTotal = exchangeMsgTotal

	systemMsgTotal, err := d.message.CountAll(repository.FindMessageParams{
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
	msgList, total, err := d.message.GetMessage(repository.FindMessageParams{
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

	uKeyMap := make(map[int64]struct{}, l+1)
	topicMap := make(map[int64]struct{}, l+1)
	commentMap := make(map[int64]struct{}, l+1)
	orderMap := make(map[int64]struct{}, l+1)
	goodsMap := make(map[int64]struct{}, l+1)
	for i, item := range msgList {
		one := &GetWebMessageResp{}
		_ = util.MapTo(item, one)
		result[i] = one
		//uKeyMap
		uKeyMap[item.SendId] = struct{}{}
		//turnMap
		if item.TurnType == 1 {
			topicMap[item.TurnId] = struct{}{}
		}
		if item.TurnType == 2 {
			commentMap[item.TurnId] = struct{}{}
		}
		if item.TurnType == 3 {
			orderMap[item.TurnId] = struct{}{}
		}
		if item.TurnType == 4 {
			goodsMap[item.TurnId] = struct{}{}
		}
	}

	//User
	var uIds []int64
	for id, _ := range uKeyMap {
		uIds = append(uIds, id)
	}

	uMap := make(map[int64]entity.ShortUser, len(uIds))
	uList := d.user.GetShortUserListBy(repository.GetUserListBy{UserIds: uIds})
	for _, uItem := range uList {
		uMap[uItem.ID] = uItem
	}

	//Turn
	tMap := d.turnTopic(topicMap)     //文章
	cMap := d.turnComment(commentMap) //评论
	oMap := d.turnOrder(orderMap)     // 订单
	gMap := d.turnGoods(goodsMap)     // 商品

	for _, item := range result {
		item.User = uMap[item.SendId]
		//文章组合
		if item.TurnType == 1 {
			item.TurnNotes = tMap[item.TurnId]
		}
		//评论组合
		if item.TurnType == 2 {
			item.TurnNotes = cMap[item.TurnId]
		}
		//商品组合
		if item.TurnType == 3 {
			item.TurnNotes = oMap[item.TurnId]
		}
		//订单组合
		if item.TurnType == 4 {
			item.TurnNotes = gMap[item.TurnId]
		}
	}

	return result, total, nil
}

func (d defaultWebMessage) SendMessage(param SendWebMessage) error {
	content := d.getTemplate(param.Key)

	if content == "" {
		return errors.New("模板不存在")
	}

	if param.TurnType == 1 {
		obj := d.topic.FindById(param.TurnId)
		content = strings.ReplaceAll(content, "{0}", obj.Title)
	}
	if param.TurnType == 2 {
		obj, _ := d.comment.FindOne(param.TurnId)
		content = strings.ReplaceAll(content, "{0}", obj.Message)
	}

	//入库
	err := d.message.SendMessage(repository.SendMessage{
		SendId:   param.SendId,
		RecId:    param.RecId,
		Type:     param.Type,
		TurnType: param.TurnType,
		TurnId:   param.TurnId,
		Message:  content,
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

func (d defaultWebMessage) getTemplate(key string) string {
	one, err := d.template.FindOne(key)
	if err != nil {
		return ""
	}
	return one.TempContent
}

func (d defaultWebMessage) turnTopic(topicMap map[int64]struct{}) map[int64]string {
	if len(topicMap) >= 1 {
		var topicIds []int64
		for id, _ := range topicMap {
			topicIds = append(topicIds, id)
		}
		topicList := d.topic.GetTopicNotes(topicIds)
		if len(topicList) >= 1 {
			tMap := make(map[int64]string, len(topicIds))
			for _, topicItem := range topicList {
				notes := ""
				n := len(topicItem.Tags)
				if n > 2 {
					n = 2
				}
				for i := 0; i < n; i++ {
					notes += topicItem.Tags[i].Name + "|"
				}
				notes += topicItem.Title
				tMap[topicItem.Id] = notes
			}
			return tMap
		}
	}

	return map[int64]string{}
}

func (d defaultWebMessage) turnComment(commentMap map[int64]struct{}) map[int64]string {
	if len(commentMap) >= 1 {
		var commentIds []int64
		for id, _ := range commentMap {
			commentIds = append(commentIds, id)
		}
		commentList := d.comment.FindListByIds(commentIds)
		if len(commentList) >= 1 {
			cMap := make(map[int64]string, len(commentIds))
			for _, commentItem := range commentList {
				cMap[commentItem.Id] = commentItem.Message
			}
			return cMap
		}
	}
	return map[int64]string{}
}

func (d defaultWebMessage) turnOrder(commentMap map[int64]struct{}) map[int64]string {
	return map[int64]string{}
}

func (d defaultWebMessage) turnGoods(commentMap map[int64]struct{}) map[int64]string {
	return map[int64]string{}
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
		topic:           repository.NewTopicModel(ctx),
		comment:         repository.NewCommentModel(ctx),
		options:         option,
	}
}
