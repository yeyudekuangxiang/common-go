package auth

type ConfigSignForm struct {
	Url string `json:"url" form:"url" binding:"required" alias:"页面地址"`
}
