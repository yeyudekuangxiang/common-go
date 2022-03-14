package activity

import "mio/controller"

type GetBocRecordListForm struct {
	ApplyStatus int8 `json:"applyStatus" form:"applyStatus" binding:"oneof=0 1 2 3 4" alias:"参与状态"`
	controller.PageFrom
}
type AddBocRecordFrom struct {
	ShareUserId int64  `json:"shareUserId" form:"shareUserId" binding:"gte=0" alias:"分享者ID"`
	Source      string `json:"source" form:"source" alias:"用户来源"`
}

type AnswerBocQuestionFrom struct {
	Right int8 `json:"right" form:"right" binding:"oneof=2 3" alias:"答题结果"`
}

type ApplySendBonus struct {
	Type string `json:"type" form:"type" binding:"oneof=apply bind boc" alias:"奖励类型"`
}
