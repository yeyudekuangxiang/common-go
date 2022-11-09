package api_types

type GetSendSmsForm struct {
	Phone      string `json:"phone" form:"phone" binding:"required" alias:"手机号"`
	Msg        string `json:"msg" form:"msg" binding:"required" alias:"短信内容"`
	TemplateId string `json:"templateId" form:"templateId" binding:"required" alias:"templateId"`
}

type SendSmsVo struct {
	ID int64 `json:"id"`
}

type GetSendZhugeForm struct {
	EventKey string `json:"eventKey" form:"eventKey" binding:"" alias:"时间名称"`
	Openid   string `json:"openid" form:"openid" binding:"" alias:"openid"`
	Data     string `json:"data" form:"data" binging:"" alias:"data"`
}

type SendZhugeVo struct {
	ID int64 `json:"id"`
}
