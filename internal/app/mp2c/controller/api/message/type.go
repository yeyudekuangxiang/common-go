package message

import "mio/internal/app/mp2c/controller"

type MessageGetTemplateIdForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=topic platform carbonpk quiz charge" alias:"模版场景"`
}

type WebMessageRequest struct {
	controller.PageFrom
	Status int    `json:"status" form:"status"`
	Types  string `json:"types" form:"types"`
}

// 已读
type HaveReadWebMessageRequest struct {
	MsgIds string `json:"msgIds" default:"" form:"msgIds"`
}

type TurnWebMessageRequest struct {
	TurnType int   `json:"turnType" form:"turnType"`
	TurnId   int64 `json:"turnId" form:"turnId"`
}
