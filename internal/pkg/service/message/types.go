package message

type sendWebMessage struct {
	SendId   int64  `json:"sendId"`
	RecId    int64  `json:"recId"`
	Key      string `json:"key"`
	RecObjId int64  `json:"recObjId"`
	Type     int    `json:"type"`
}
