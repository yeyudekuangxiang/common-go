package activity

import (
	"mio/model/entity"
	activityM "mio/model/entity/activity"
)

type GetRecordPageListParam struct {
	UserId      int64 `json:"userId"`
	ApplyStatus int8  `json:"applyStatus"`
	Offset      int   `json:"offset"`
	Limit       int   `json:"limit"`
}
type AddApplyRecordParam struct {
	UserId      int64 `json:"userId"`
	ShareUserId int64 `json:"shareUserId"` //分享者用户id
}

type BocRecordDetail struct {
	activityM.BocRecord
	User entity.ShortUser `json:"user"`
}

type CreateBocShareBonusRecordParam struct {
	UserId int64
	Value  int64 //金额 单位分
}
