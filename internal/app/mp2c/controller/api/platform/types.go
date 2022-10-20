package platform

type AllReceiveRequest struct {
	PlatformKey string `json:"platformKey" form:"platformKey" alias:"platformKey" binding:"required"`
}
