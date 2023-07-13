package message

import (
	"github.com/medivhzhan/weapp/v3/subscribemessage"
	"strconv"
)

type IMiniSubTemplate interface {
	ToData() map[string]subscribemessage.SendValue
	TemplateId() string
	SendMixCount() float64
	Page() string
	IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU()
	TemplateName() string
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
	return MessageTemplateIds.ChangePoint
}

func (m MiniChangePointTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.ChangePoint
}

func (m MiniChangePointTemplate) Page() string {
	return MessageJumpUrls.ChangePoint
}

func (m MiniChangePointTemplate) TemplateName() string {
	return MessageTemplateName.ChangePoint
}

func (m MiniChangePointTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}

/**
详细内容
打卡主题  {{thing1.DATA}}
打卡名称   {{thing2.DATA}}
已打卡天数  {{thing7.DATA}}
提醒内容  {{thing5.DATA}}
备注  {{thing4.DATA}}
*/

//MiniClockRemindTemplate 打卡提醒

type MiniClockRemindTemplate struct {
	Title   string `json:"title"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	Content string `json:"content"`
	Tip     string `json:"tip"`
}

func (m MiniClockRemindTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"thing1": {Value: m.Title},
		"thing2": {Value: m.Name},
		"thing7": {Value: m.Date},
		"thing5": {Value: m.Content},
		"thing4": {Value: m.Tip},
	}
}

func (m MiniClockRemindTemplate) TemplateId() string {
	return MessageTemplateIds.PunchClockRemind
}

func (m MiniClockRemindTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.PunchClockRemind
}

func (m MiniClockRemindTemplate) Page() string {
	return MessageJumpUrls.PunchClockRemind
}

func (m MiniClockRemindTemplate) TemplateName() string {
	return MessageTemplateName.PunchClockRemind
}

func (m MiniClockRemindTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
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
	return MessageTemplateIds.OrderDeliver
}

func (m MiniOrderDeliverTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.OrderDeliver
}

func (m MiniOrderDeliverTemplate) Page() string {
	return MessageJumpUrls.OrderDeliver
}

func (m MiniOrderDeliverTemplate) TemplateName() string {
	return MessageTemplateName.OrderDeliver
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
	return MessageTemplateIds.SignRemind
}

func (m MiniSignRemindTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.SignRemind
}

func (m MiniSignRemindTemplate) Page() string {
	return MessageJumpUrls.SignRemind
}

func (m MiniSignRemindTemplate) TemplateName() string {
	return MessageTemplateName.SignRemind
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
	return MessageTemplateIds.TopicPass
}

func (m MiniTopicPassTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.TopicPass
}

func (m MiniTopicPassTemplate) Page() string {
	return MessageJumpUrls.TopicPass
}

func (m MiniTopicPassTemplate) TemplateName() string {
	return MessageTemplateName.TopicPass
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
	return MessageTemplateIds.TopicCarefullyChosen
}

func (m MiniTopicCarefullyChosenTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.TopicCarefullyChosen
}

func (m MiniTopicCarefullyChosenTemplate) Page() string {
	return MessageJumpUrls.TopicCarefullyChosen
}

func (m MiniTopicCarefullyChosenTemplate) TemplateName() string {
	return MessageTemplateName.TopicCarefullyChosen
}

func (m MiniTopicCarefullyChosenTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}

//MiniQuizRemindTemplate 每日答题闯关提醒
type MiniQuizRemindTemplate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Remark  string `json:"remark"`
}

func (m MiniQuizRemindTemplate) ToData() map[string]subscribemessage.SendValue {
	return map[string]subscribemessage.SendValue{
		"thing1": {Value: m.Title},
		"thing2": {Value: m.Content},
		"thing3": {Value: m.Remark},
	}
}

func (m MiniQuizRemindTemplate) TemplateId() string {
	return MessageTemplateIds.QuizRemind
}

func (m MiniQuizRemindTemplate) SendMixCount() float64 {
	return MessageSendMixCounts.QuizRemind
}

func (m MiniQuizRemindTemplate) Page() string {
	return MessageJumpUrls.QuizRemind
}
func (m MiniQuizRemindTemplate) TemplateName() string {
	return MessageTemplateName.QuizRemind
}
func (m MiniQuizRemindTemplate) IDIUDGIUGISIAHSUIAHUISAHUIAGUISGIU() {
	return
}
