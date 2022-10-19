package jhx

type jhxCommonResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Time string                 `json:"time"`
	Data map[string]interface{} `json:"data"`
}

type jhxCommonRequest struct {
	Version string `json:"version"`
	Appid   string `json:"appid"`
	Nonce   int    `json:"nonce"`
}

type jhxTicketNotifyRequest struct {
	jhxCommonRequest
	TicketNo string `json:"ticket_no" form:"ticket_no"`
	Status   string `json:"status" form:"status"`
	UsedTime string `json:"used_time" form:"used_time"`
}

type jhxTicketCreateResponse struct {
	QrCodeStr  string `json:"qrcodestr" form:"qrcodestr"`
	ExpireTime string `json:"expireTime" from:"expireTime"`
}

type jhxTicketStatusResponse struct {
	TicketNo string `json:"ticket_no" form:"ticket_no"`
	Status   string `json:"status" form:"status"`
	UsedTime string `json:"used_time" form:"used_time"`
}

// request
type Collect struct {
	MemberId    string `json:"memberId" from:"memberId"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	OpenId      string `json:"openId" form:"openId"`
	PrePointId  string `json:"prePointId" form:"prePointId"`
}

type GetPreCollect struct {
	MemberId    string `json:"memberId" from:"memberId"`
	OpenId      string `json:"openId" form:"openId"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
}

// callback
type TicketNotify struct {
	Tradeno  string `json:"tradeno" form:"tradeno" binding:"required"`
	Status   string `json:"status" form:"status" binding:"required"`
	UsedTime string `json:"used_time" form:"used_time" binding:"required"`
	Sign     string `json:"sign" form:"sign" binding:"required"`
}
