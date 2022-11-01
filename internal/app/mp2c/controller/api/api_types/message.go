package api_types

import "mio/internal/app/mp2c/controller"

type MessageGetTemplateIdForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=topic platform" alias:"模版场景"`
}

type WebMessageRequest struct {
	controller.PageFrom
	Status int `json:"status" form:"status" binding:"required"`
	Type   int `json:"type" form:"type" binding:"required"`
}

// 已读
type HaveReadWebMessageRequest struct {
	MsgId  int64   `json:"msgId" form:"msgId"`
	RecId  int64   `json:"recId" form:"recId" binding:"required"`
	MsgIds []int64 `json:"msgIds" form:"msgIds"`
}
