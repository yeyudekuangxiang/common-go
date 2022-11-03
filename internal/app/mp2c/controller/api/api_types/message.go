package api_types

import "mio/internal/app/mp2c/controller"

type MessageGetTemplateIdForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=topic platform" alias:"模版场景"`
}

type WebMessageRequest struct {
	controller.PageFrom
	Status int    `json:"status" default:"1" form:"status" binding:"required"`
	Types  string `json:"types" form:"types"`
}

// 已读
type HaveReadWebMessageRequest struct {
	MsgIds string `json:"msgIds" form:"msgIds"`
}
