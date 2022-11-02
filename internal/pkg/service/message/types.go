package message

type SendWebMessage struct {
	SendId   int64  `json:"sendId"`
	RecId    int64  `json:"recId"`
	Key      string `json:"key"`
	RecObjId int64  `json:"recObjId"`
	Type     int    `json:"type" default:"1"`
}

type SetHaveReadMessage struct {
	MsgId  int64   `json:"msgId" form:"msgId"`
	MsgIds []int64 `json:"msgIds" form:"msgIds"`
	RecId  int64   `json:"recId" form:"recId"`
}

type GetWebMessage struct {
	UserId int64 `json:"userId"`
	Status int   `json:"status"`
	Type   int   `json:"type"`
	Types  []int `json:"types"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type GetWebMessageCount struct {
	RecId int64 `json:"recId"`
}

type GetWebMessageCountResp struct {
	Total            int64 `json:"total"`
	ExchangeMsgTotal int64 `json:"exchangeMsgTotal"`
	SystemMsgTotal   int64 `json:"systemMsgTotal"`
}
