package message

type SendWebMessage struct {
	SendId   int64  `json:"sendId"`
	RecId    int64  `json:"recId"`
	Key      string `json:"key"`
	RecObjId int64  `json:"recObjId"`
	Type     int    `json:"type"`
}

type SetHaveReadMessage struct {
	MsgId  int64   `json:"msgId" form:"msgId"`
	MsgIds []int64 `json:"msgIds" form:"msgIds"`
	RecId  int64   `json:"recId" form:"recId"`
}
