package activity

import (
	"mio/model"
	"mio/model/entity"
	activityM "mio/model/entity/activity"
)

type GetRecordPageListParam struct {
	UserId                  int64 `json:"userId"`
	ApplyStatus             int8  `json:"applyStatus"`
	ShareUserBocBonusStatus int8  `json:"shareUserBocBonusStatus"`
	Offset                  int   `json:"offset"`
	Limit                   int   `json:"limit"`
}
type AddApplyRecordParam struct {
	UserId      int64  `json:"userId"`
	ShareUserId int64  `json:"shareUserId"` //分享者用户id
	Source      string `json:"source"`
}

type BocRecordDetail struct {
	activityM.BocRecord
	CreatedAtDate model.Date       `json:"createdAtDate"`
	UpdatedAtDate model.Date       `json:"updatedAtDate"`
	User          entity.ShortUser `json:"user"`
}

type CreateBocShareBonusRecordParam struct {
	UserId int64
	//金额 单位分
	Value int64
	Type  activityM.BocShareBonusType
	Info  string
}
type AnswerGMQuestionParam struct {
	UserId  int64
	Title   string
	Answer  string
	IsRight bool
}
