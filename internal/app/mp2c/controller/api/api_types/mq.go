package api_types

type GetSendSmsForm struct {
	Phone       string `json:"phone" form:"phone" binding:"required" alias:"手机号"`
	Msg         string `json:"msg" form:"msg" binding:"required" alias:"短信内容"`
	TemplateKey string `json:"templateKey" form:"templateKey" binding:"required" alias:"templateKey"`
}
type GetSendYzmSmsForm struct {
	Phone string `json:"phone" form:"phone" binding:"required" alias:"手机号"`
	Code  string `json:"code" form:"code" binding:"required" alias:"验证码"`
}

type GetSendZhugeForm struct {
	EventKey string `json:"eventKey" form:"eventKey" binding:"required" alias:"事件名称"`
	Openid   string `json:"openid" form:"openid" binding:"required" alias:"openid"`
	Data     string `json:"data" form:"data" binging:"" alias:"data"`
}
