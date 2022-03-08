package activity

import "mio/controller"

type GetRecordListForm struct {
	ApplyStatus int8 `json:"applyStatus" form:"applyStatus" binding:"oneof=0 1 2 3 4" alias:"参与状态"`
	controller.PageFrom
}
