package event

type GetEventFullDetailForm struct {
	EventId string `json:"eventId" form:"eventId" binding:"required" alias:"项目编号"`
}
type GetEventListForm struct {
	EventCategoryId string `json:"eventCategoryId" form:"eventCategoryId" binding:"required" alias:"项目类型"`
}
