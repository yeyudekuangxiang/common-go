package open

type platformForm struct {
	PlatformKey string `json:"platformKey" form:"platformKey"`
	MemberId    string `json:"memberId,omitempty" form:"memberId"`
	Method      string `json:"method,omitempty" form:"method"`
	Mobile      string `json:"mobile,omitempty" form:"mobile"`
	Sign        string `json:"sign,omitempty" form:"sign"`
}

type bindPlatform struct {
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	MemberId    string `json:"memberId" form:"memberId" binding:"required"`
	IsNewUser   string `json:"isNewUser" form:"isNewUser"`
}

//金华行 核销参数
type jhxTicketNotifyRequest struct {
	Tradeno  string `json:"tradeno" form:"tradeno" binding:"required"`
	Status   string `json:"status" form:"status" binding:"required"`
	UsedTime string `json:"used_time" form:"used_time" binding:"required"`
	Sign     string `json:"sign" form:"sign" binding:"required"`
}

type jhxTicketStatusRequest struct {
	Tradeno string `json:"tradeno" form:"tradeno" binding:"required"`
}

type jhxTicketCreateRequest struct {
	Tradeno   string `json:"tradeno" form:"tradeno" binding:"required"`
	TypeId    int64  `json:"typeId" form:"typeId" binding:"required"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
}

type jhxPreCollectRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	Amount      string `json:"amount" form:"amount" binding:"required"`
	//PrePointId  string `json:"prePointId,omitempty" form:"prePointId"`
	Sign string `json:"sign" form:"sign" binding:"required"`
}

type jhxGetPreCollectRequest struct {
	MemberId    string `json:"memberId" from:"memberId"`
	OpenId      string `json:"openId" form:"openId"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	Sign        string `json:"sign" form:"sign" binding:"required"`
}

type jhxCollectRequest struct {
	MemberId    string `json:"memberId" from:"memberId"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	OpenId      string `json:"openId" form:"openId"`
	PrePointId  string `json:"prePointId" form:"prePointId" binding:"required"`
	Sign        string `json:"sign" form:"sign" binding:"required"`
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

type allReceiveRequest struct {
	PlatformKey string `json:"platformKey" form:"platformKey" alias:"platformKey" binding:"required"`
}

type PrePointListRequestV2 struct {
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
}

type setPrePointRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	Point       int64  `json:"point" form:"point" binding:"required"`
	TradeNo     string `json:"tradeNo" form:"tradeNo"  binding:"required"`
	Mobile      string `json:"mobile,omitempty" form:"mobile"`
	Sign        string `json:"sign" form:"sign" binding:"required"`
}

type prePointListRequest struct {
	MemberId    string `json:"memberId" from:"memberId" binding:"required"`
	PlatformKey string `json:"platformKey" form:"platformKey" binding:"required"`
	Sign        string `json:"sign" form:"sign" binding:"required"`
}

type checkMsg struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type checkMedia struct {
	MediaUrl string `json:"mediaUrl" form:"mediaUrl" binding:"required"`
}

type recycleReq struct {
	Ch           string `json:"ch" form:"ch" binding:"required"`
	OrderNo      string `json:"orderNo" form:"orderNo" binding:"required"`
	MemberId     string `json:"memberId" form:"memberId" binding:"required"`
	Name         string `json:"name" form:"name" binding:"required"`
	Category     string `json:"category" form:"category" binding:"required"`
	Number       string `json:"number" form:"number" binding:"required"`
	CreateTime   string `json:"createTime"`
	CompleteTime string `json:"completeTime"`
	T            string `json:"t"`
	Sign         string `json:"sign" form:"sign" binding:"required"`
}

type YkcReq struct {
	ExternalUserId      string  `json:"externalUserId" binding:"required"`
	TradeSeq            string  `json:"tradeSeq" binding:"required"`
	GunName             string  `json:"gunName" binding:"required"`
	StationName         string  `json:"stationName" binding:"required"`
	ChargedActualAmount float64 `json:"chargedActualAmount"`
	ServiceActualAmount float64 `json:"serviceActualAmount"`
	ActualTotalAmount   float64 `json:"actualTotalAmount"`
	ChargedPower        float64 `json:"chargedPower" binding:"required"`
	StartTime           string  `json:"startTime" binding:"required"`
	EndTime             string  `json:"endTime"`
	TradeStatus         string  `json:"tradeStatus" binding:"required"`
	PaymentWallet       string  `json:"paymentWallet" binding:"required"`
}

type incPointForActivityParams struct {
	Point        int64  `json:"point"`
	OpenId       string `json:"openId"`
	AdditionInfo string `json:"additionInfo"`
}
