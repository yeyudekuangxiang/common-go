package api_types

type SubmitOrderForEventForm struct {
	EventId string `json:"eventId" form:"eventId" binding:"required" alias:"项目编号"`
}

type SubmitOrderForEventGDForm struct {
	EventId        string `json:"eventId" form:"eventId" binding:"required" alias:"项目编号"`
	WxServerOpenId string `json:"wxServerOpenId" form:"wxServerOpenId" alias:"微信服务号"`
}
