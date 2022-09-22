package open

type platformForm struct {
	PlatformKey string `json:"platformKey" form:"platformKey"`
	MemberId    string `json:"memberId,omitempty" form:"memberId"`
	Method      string `json:"method,omitempty" form:"method"`
	Mobile      string `json:"mobile,omitempty" form:"mobile"`
	Sign        string `json:"sign,omitempty" form:"sign"`
}

//金华行 核销参数
type jhxCommonForm struct {
	Version string `json:"version,omitempty"`
	Appid   string `json:"appid,omitempty"`
	Nonce   int    `json:"nonce,omitempty"`
	Sign    string `json:"sign,omitempty"`
}
type jhxUseCodeForm struct {
	jhxCommonForm
	TicketNo string `json:"ticket_no" form:"ticket_no" binding:"required"`
	Status   string `json:"status" form:"status" binding:"required"`
	UsedTime string `json:"used_time" form:"used_time"`
}
