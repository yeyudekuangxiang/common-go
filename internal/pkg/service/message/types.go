package message

import (
	"mio/internal/pkg/model/entity"
	"time"
)

//消息模板
const (
	//二手市场
	SecondHandGetPoint                 = "secondHand-get_point"
	SecondHandGetPointCommodityPublish = "secondHand-commodity_publish"
	SecondHandGetPointLikeComment      = "secondHand-like_comment"
	SecondHandGetPointReplyComment     = "secondHand-reply_comment"
	SecondHandLikeCommodity            = "secondHand-like_commodity"
	SecondHandCommentCommodity         = "secondHand-comment_commodity"
	SecondHandDeliverTimeOut           = "secondHand-deliver-time-out"
	SecondHandOrderPayed               = "secondHand-order-payed"
	SecondHandOrderCancel              = "secondHand-order-cancel"
	SecondHandReceiveTimeOut           = "secondHand-receive-time-out"
	SecondHandOrderReceive             = "secondHand-order-receive"
	SecondHandOrderDeliver             = "secondHand-order-deliver"
	SecondHandOrderAward               = "secondHand-order-award"
	SecondHandAuditReject              = "secondHand-audit_reject"
	SecondHandSellerDeliverTimeOut     = "secondHand-seller-deliver-time-out"
	FailGoods                          = "fail_goods"

	//社区
	BanCommunity   = "ban_community"
	ListTopic      = "like_topic"
	PushTopic      = "push_topic"
	PushTopicV2    = "push_topic_v2"
	DownComment    = "down_comment"
	FailTopic      = "fail_topic"
	Wechat         = "wechat"
	LikeComment    = "like_comment"
	ReplyTopic     = "reply_topic"
	TopTopic       = "top_topic"
	TopTopicV2     = "top_topic_v2"
	DownTopic      = "down_topic"
	ReplyComment   = "reply_comment"
	EssenceTopic   = "essence_topic"
	EssenceTopicV2 = "essence_topic_v2"

	//短信
	SmsSecondHandOrderPayed           = "sms-secondHand-order-payed"
	SmsSecondHandOrderCancel          = "sms-secondHand-order-cancel"
	SmsSecondHandDeliverTimeOut       = "sms-secondHand-deliver-time-out"
	SmsSecondHandOrderDeliver         = "sms-secondHand-order-deliver"
	SmsSecondHandOrderReceive         = "sms-secondHand-order-receive"
	SmsSecondHandSellerDeliverTimeOut = "sms-secondHand-seller-deliver-time-out"
	SmsSuiShenXingLuckydraw           = "sms-suishenxing-luckydraw"
	SmsPointTimeExpireQuickly         = "sms-point-time-expire-quickly"
	SmsPointTimeExpire                = "sms-point-time-expire"
	SmsActivityCancel                 = "sms-activity-cancel"
)

const (
	MsgTypeNone    = iota
	MsgTypeLike    //点赞
	MsgTypeComment //评论||回复
	MsgTypeSystem  //系统消息
)

const (
	MsgTurnTypeNone           = iota
	MsgTurnTypeArticle        //跳转文章
	MsgTurnTypeArticleComment //跳转文章评论
	MsgTurnTypeOrder          //跳转订单
	MsgTurnTypeGoods          //跳转商品
	MsgTurnTypeGoodsComment   //跳转商品评论
)

type SendWebMessage struct {
	SendId       int64  `json:"sendId"`
	RecId        int64  `json:"recId"`
	Key          string `json:"key"`
	Type         int    `json:"type"`     // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社 8 商品 9订单
	TurnType     int    `json:"turnType"` // object_turn_type 1文章 2评论 3订单 4商品
	TurnId       int64  `json:"turnId"`   // object_turn_id   要跳转的object的id
	MessageNotes string `json:"messageNotes"`
}

type SetHaveReadMessage struct {
	MsgIds []string `json:"msgIds" form:"msgIds"`
	RecId  int64    `json:"recId" form:"recId"`
}

type GetWebMessage struct {
	UserId int64    `json:"userId"`
	Status int      `json:"status"`
	Type   int      `json:"type"`
	Types  []string `json:"types"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

type GetWebMessageResp struct {
	//message
	Id             int64     `json:"id"`
	MessageContent string    `json:"messageContent"`
	MessageNotes   string    `json:"messageNotes"`
	Type           int       `json:"type"`   // 1点赞 2评论 3回复 4发布 5精选 6违规 7合作社
	Status         int       `json:"status"` // 1未读 2已读
	CreatedAt      time.Time `json:"createdAt"`
	//turn obj
	TurnType int    `json:"turnType"` // 1文章 2评论 3订单 4商品
	TurnId   string `json:"turnId"`
	//user
	SendId int64            `json:"sendId"`
	User   entity.ShortUser `json:"user"`
}

type GetWebMessageCount struct {
	RecId int64 `json:"recId"`
}

type GetWebMessageCountResp struct {
	Total               int64 `json:"total"`
	InteractiveMsgTotal int64 `json:"exchangeMsgTotal"`
	SystemMsgTotal      int64 `json:"systemMsgTotal"`
}

type GetMessageTemplate struct {
	Id          int64     `json:"id"`
	Key         string    `json:"key"`
	Type        int       `json:"type"`
	TempContent string    `json:"tempContent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type jsonObj struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	DelReason string `json:"delReason"`
}
