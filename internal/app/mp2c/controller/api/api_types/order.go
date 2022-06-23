package api_types

type SubmitOrderForEventForm struct {
	EventId string `json:"eventId" form:"eventId" binding:"required" alias:"项目编号"`
}
