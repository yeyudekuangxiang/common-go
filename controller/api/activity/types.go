package activity

import "mio/controller"

type GetBocApplyRecordListForm struct {
	ApplyStatus int8 `json:"applyStatus" form:"applyStatus" binding:"oneof=0 1 2 3 4" alias:"参与状态"`
	controller.PageFrom
}
type AddBocApplyRecordFrom struct {
	ShareUserId int64 `json:"shareUserId" form:"shareUserId" binding:"gte=0" alias:"分享者ID"`
}
