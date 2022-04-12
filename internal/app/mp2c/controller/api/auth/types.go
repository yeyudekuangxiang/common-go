package auth

type oa struct {
	Platform string `json:"platform" form:"platform" binding:"required,oneof=mio-srv-oa mio-sub-oa" alias:"platform"`
}
type ConfigSignForm struct {
	oa
	Url string `json:"url" form:"url" binding:"required" alias:"页面地址"`
}
type AutoLoginForm struct {
	oa
	RedirectUri string `json:"redirectUri" form:"redirectUri" binding:"required" alias:"redirectUri"`
	State       string `json:"state" form:"state"`
}
type AutoLoginCallbackForm struct {
	oa
	Code  string `json:"code" form:"code" binding:"required" alias:"code"`
	State string `json:"state" form:"code" binding:"required" alias:"code"`
}
type OaAuthForm struct {
	oa
	Code string `json:"code" form:"code" binding:"required" alias:"code"`
}
