package jhx

type jhxCommonResponse struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Time string                 `json:"time"`
	Data map[string]interface{} `json:"data"`
}

type jhxTicketCreateResponse struct {
	QrCodeStr string `json:"qrcodestr" form:"qrcodestr"`
}

type jhxTicketStatusResponse struct {
	TicketNo string `json:"ticket_no" form:"ticket_no"`
	Status   string `json:"status" form:"status"`
	UsedTime string `json:"used_time" form:"used_time"`
}
