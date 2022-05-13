package auth

type oa struct {
	Platform string `json:"platform" form:"platform" binding:"oneof=mio-srv-oa mio-sub-oa" alias:"platform"`
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
	Code  string `json:"code" form:"code" binding:"required" alias:"code"`
	State string `json:"state" form:"state" binding:"required" alias:"state"`
}
type OaAuthForm struct {
	oa
	Code string `json:"code" form:"code" binding:"required" alias:"code"`
}
type WeappAuthForm struct {
	Code            string `json:"code" form:"code" binding:"required" alias:"code"`
	PartnershipWith string `json:"partnershipWith" form:"partnershipWith" alias:"partnershipWith"`
	InvitedBy       string `json:"invitedBy" form:"invitedBy" alias:"invitedBy"`
}
