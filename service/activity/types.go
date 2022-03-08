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

type BocApplyRecordDetail struct {
	activityM.BocApplyRecord
	User entity.ShortUser `json:"user"`
}
