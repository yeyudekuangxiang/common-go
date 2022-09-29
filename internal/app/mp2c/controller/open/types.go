package open

type platformForm struct {
	PlatformKey string `json:"platformKey" form:"platformKey"`
	MemberId    string `json:"memberId,omitempty" form:"memberId"`
	Method      string `json:"method,omitempty" form:"method"`
	Mobile      string `json:"mobile,omitempty" form:"mobile"`
	Sign        string `json:"sign,omitempty" form:"sign"`
}

//金华行 核销参数
type jhxCommonRequest struct {
	Version string `json:"version,omitempty"`
	Appid   string `json:"appid,omitempty"`
	Nonce   int    `json:"nonce,omitempty"`
	Sign    string `json:"sign,omitempty"`
}

type jhxTicketNotifyRequest struct {
	jhxCommonRequest
	TicketNo string `json:"ticket_no" form:"ticket_no" binding:"required"`
	Status   string `json:"status" form:"status" binding:"required"`
	UsedTime string `json:"used_time" form:"used_time"`
}

type jhxTicketStatusRequest struct {
	Tradeno string `json:"tradeno" form:"tradeno" binding:"required"`
}

type jhxPreCollectRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	Amount      string `json:"amount" form:"amount" binding:"required"`
	//PrePointId  string `json:"prePointId,omitempty" form:"prePointId"`
	Sign string `json:"sign" form:"sign" binding:"required"`
}

type jhxGetPreCollectRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	//Amount      string `json:"amount,omitempty" form:"amount"`
	//PrePointId string `json:"prePointId,omitempty" form:"prePointId"`
	Sign string `json:"sign" form:"sign" binding:"required"`
}

type jhxCollectRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	//Amount      string `json:"amount,omitempty" form:"amount"`
	PrePointId string `json:"prePointId" form:"prePointId"`
	Sign       string `json:"sign" form:"sign" binding:"required"`
}
type jhxMyOrderRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	//Amount      string `json:"amount,omitempty" form:"amount"`
	//PrePointId  string `json:"prePointId,omitempty" form:"prePointId"`
	Sign string `json:"sign" form:"sign" binding:"required"`
}

type jhxMyAccountRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	//Amount      string `json:"amount,omitempty" form:"amount"`
	//PrePointId  string `json:"prePointId,omitempty" form:"prePointId"`
	Sign string `json:"sign" form:"sign" binding:"required"`
}

type jhxMyCrRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	//Amount      string `json:"amount,omitempty" form:"amount"`
	//PrePointId  string `json:"prePointId,omitempty" form:"prePointId"`
	Sign string `json:"sign" form:"sign"`
}
