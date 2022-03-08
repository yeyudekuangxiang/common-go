package activity

import (
	"mio/model/entity"
	activityM "mio/model/entity/activity"
)

type GetRecordPageListParam struct {
	UserId            int64 `json:"userId"`
	ApplyRecordStatus int8  `json:"applyRecordStatus"`
	Offset            int
	Limit             int
}

type BocApplyRecordDetail struct {
	activityM.BocApplyRecord
	ShareUser entity.User `json:"shareUser"`
}
