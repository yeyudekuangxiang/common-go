package message

import (
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"mio/config"
	"strconv"
)

type IMiniSubTemplate interface {
	ToData() map[string]subscribemessage.SendValue
	TemplateId() string
	IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU()
}

/*
获取积分 {{number1.DATA}}
获取来源 {{thing5.DATA}}
获取时间 {{time2.DATA}}
累积积分 {{number3.DATA}}
*/

//MiniChangePointTemplate 积分变更
type MiniChangePointTemplate struct {
	Point    int64  `json:"point"`
	Source   string `json:"source"`
	Time     string `json:"time"`
	AllPoint int64  `json:"tip"`
}

func (m MiniChangePointTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"number1": {Value: strconv.FormatInt(m.Point, 10)}, //strconv.Itoa(m.Money)strconv.Atoi(m.Point)
		"thing5":  {Value: m.Source},
		"time2":   {Value: m.Time},
		"number3": {Value: strconv.FormatInt(m.AllPoint, 10)},
	}
}

func (m MiniChangePointTemplate) TemplateId() string {
	return config.MessageTemplateIds.ChangePoint
}

func (m MiniChangePointTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}

/**
"订单号 {{character_string1.DATA}}
物流单号 {{character_string3.DATA}}
物流公司 {{thing4.DATA}}
商品名称 {{thing5.DATA}}
温馨提示 {{thing6.DATA}}"
*/

//MiniOrderDeliverTemplate 订单发货提醒
type MiniOrderDeliverTemplate struct {
	OrderNo      string `json:"orderNo"`
	TrackNo      string `json:"trackNo"`
	TrackCompany string `json:"trackCompany"`
	GoodName     string `json:"goodName"`
	Tip          string `json:"tip"`
}

func (m MiniOrderDeliverTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"character_string1": {Value: m.OrderNo}, //strconv.Itoa(m.Money)strconv.Atoi(m.Point)
		"character_string3": {Value: m.TrackNo},
		"thing4":            {Value: m.TrackCompany},
		"thing5":            {Value: m.GoodName},
		"thing6":            {Value: m.Tip},
	}
}

func (m MiniOrderDeliverTemplate) TemplateId() string {
	return config.MessageTemplateIds.OrderDeliver
}

func (m MiniOrderDeliverTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}

/**
活动名称 {{thing1.DATA}}
温馨提醒 {{thing9.DATA}}
*/

//MiniSignRemindTemplate 签到提醒
type MiniSignRemindTemplate struct {
	ActivityName string `json:"activityName"`
	Tip          string `json:"tip"`
}

func (m MiniSignRemindTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"thing1": {Value: m.ActivityName},
		"thing9": {Value: m.Tip},
	}
}

func (m MiniSignRemindTemplate) TemplateId() string {
	return config.MessageTemplateIds.SignRemind
}

func (m MiniSignRemindTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}

/**
审核结果  {{thing1.DATA}}
帖子主题  {{thing2.DATA}}
通过时间  {{date4.DATA}}
备注     {{thing5.DATA}}
*/

//MiniTopicPassTemplate 帖子审核通过提醒
type MiniTopicPassTemplate struct {
	AuditResult string `json:"auditResult"`
	TopicTitle  string `json:"topicTitle"`
	Time        string `json:"time"`
	Tip         string `json:"tip"`
}

func (m MiniTopicPassTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"thing1": {Value: m.AuditResult},
		"thing2": {Value: m.TopicTitle},
		"date4":  {Value: m.Time},
		"thing5": {Value: m.Tip},
	}
}

func (m MiniTopicPassTemplate) TemplateId() string {
	return config.MessageTemplateIds.TopicPass
}

func (m MiniTopicPassTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}

/**
帖子标题  {{thing1.DATA}}
时间         {{time2.DATA}}
获取积分  {{number3.DATA}}
备注         {{thing4.DATA}}
*/

//MiniTopicCarefullyChosenTemplate 帖子被精选通知
type MiniTopicCarefullyChosenTemplate struct {
	TopicTitle string `json:"topicTitle"`
	Time       string `json:"time"`
	Point      string `json:"point"`
	Tip        string `json:"tip"`
}

func (m MiniTopicCarefullyChosenTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"thing1":  {Value: m.TopicTitle},
		"time2":   {Value: m.Time},
		"number3": {Value: m.Point},
		"thing4":  {Value: m.Tip},
	}
}

func (m MiniTopicCarefullyChosenTemplate) TemplateId() string {
	return config.MessageTemplateIds.TopicCarefullyChosen
}

func (m MiniTopicCarefullyChosenTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}
