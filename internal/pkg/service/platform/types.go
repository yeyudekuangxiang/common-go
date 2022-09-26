package platform

type FmySignParams struct {
	AppId          string         `json:"app_id" form:"app_id" binding:"required"`
	NotificationAt string         `json:"notification_at" form:"notification_at" binding:"required"`
	Data           recycleFmyData `json:"data" form:"data"`
	Sign           string         `json:"sign,omitempty"`
}

type recycleFmyData struct {
	OrderSn          string `json:"order_sn,omitempty" binding:"required"`
	Status           string `json:"status,omitempty" binding:"required"`
	Weight           string `json:"weight,omitempty"`
	Reason           string `json:"reason,omitempty"`
	CourierRealName  string `json:"courier_real_name,omitempty"`
	CourierPhone     string `json:"courier_phone,omitempty"`
	CourierJobNumber string `json:"courier_job_number,omitempty"`
	Waybill          string `json:"waybill,omitempty"`
	Phone            string `json:"phone,omitempty"`
}

type oolaSignParams struct {
	Type                int    `json:"type" form:"type" binding:"required"`                     //业务类型 1：回首订单成功
	OrderNo             string `json:"orderNo" form:"orderNo" binding:"required"`               //订单号，同类型同订单视为重复订单
	Name                string `json:"name" form:"name" binding:"required"`                     //type = 1，回收物品名称
	ProductCategoryName string `json:"productCategoryName" form:"productCategoryName"`          //物品所属分类名称
	Qua                 string `json:"qua" form:"qua"`                                          //用户下单时的数量&重量
	Unit                string `json:"unit" form:"unit"`                                        //与下单数量&重量关联的计量单位 如：公斤，个 等
	OolaUserId          int    `json:"oolaUserId" form:"oolaUserId" binding:"required"`         //噢啦平台用户id
	ClientId            string `json:"clientId" form:"clientId" binding:"required"`             //lvmiao用户id
	CreateTime          string `json:"createTime" form:"createTime" binding:"required"`         //订单创建时间
	CompletionTime      string `json:"completionTime" form:"completionTime" binding:"required"` //订单完成时间
	//Sign                string `json:"sign" form:"sign" binding:"required"`                     //加密串
}

type ccRingReqParams struct {
	MemberId       string  `json:"memberId"`
	DegreeOfCharge float64 `json:"degreeOfCharge"`
}

type starChargeResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data string `json:"Data"`
	Sig  string `json:"Sig"`
}

type starChargeAccessResult struct {
	OperatorID         string `json:"operatorID,omitempty"`
	SucStat            int    `json:"sucStat,omitempty"`
	AccessToken        string `json:"accessToken,omitempty"`
	TokenAvailableTime int    `json:"tokenAvailableTime,omitempty"`
	FailReason         int    `json:"failReason,omitempty"`
}

type starChargeProvideResult struct {
	SuccStat      int    `json:"succStat,omitempty"`
	FailReason    int    `json:"failReason,omitempty"`
	FailReasonMsg string `json:"failReasonMsg,omitempty"`
	CouponCode    string `json:"couponCode,omitempty"`
}

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