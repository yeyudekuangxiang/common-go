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
