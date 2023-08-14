package growthsystemmsg

type GrowthSystemParam struct {
	TaskType    string `json:"taskType"`    //大类 非必填
	TaskSubType string `json:"taskSubType"` //具体类型
	UserId      string `json:"userId"`
	TaskValue   int64  `json:"taskValue"`
}

type GrowthSystemReq struct {
	TaskType    string `json:"taskType"`
	TaskSubType string `json:"taskSubType"`
	UserId      string `json:"userId"`
	TaskValue   string `json:"taskValue"`
	MessageId   string `json:"messageId"`
	Time        string `json:"time"`
}
