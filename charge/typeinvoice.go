package charge

type NotificationMspPaymentInfoParam struct {
	StartChargeSeq       string  //互联互通充电订单号
	UserPaidAmount       float64 //用户实付金额
	SubsidyElecAmount    float64 //用户补贴优惠电费金额
	SubsidyServiceAmount float64 //用户补贴优惠服务费金额
	RoyaltyServiceAmount float64 //互联互通渠道方服务费抽成金额
	ChannelAmount        float64 //用户支付通道费金额
}

type NotificationMspPaymentInfoResult struct {
	Status int `json:"Status"`
}

type InvoiceApply struct {
	OrderType    int
	OutInvoiceId string
	InvoiceOrders
}

type InvoiceOrders struct {
	StartChargeSeq string
	ElecMoney      float64
	SeviceMoney    float64
	TotalMoney     float64
}

type InvoiceApplyParam struct {
	OrderType       int
	OutInvoiceId    string
	InvoiceOrders   []InvoiceOrders
	BusinessType    int
	InvoiceTitle    string
	InvoiceType     int
	PayerRegisterNo string
	Remark          string
	PayerEmail      string
	ReceiverName    string
	ReceiverAddress string
	ReceiverPhone   string
	CompanyPhone    string
	CompanyAddress  string
	BankName        string
	Account         string
}

type InvoiceApplyResult struct {
	BizCode string
	BizMsg  string
	Batches []Batches
}

type Batches struct {
	BatchNo        string
	StartChargeSeq string
}

type InvoiceInfoParam struct {
	OutInvoiceId string
	BatchNo      string
	SubBatchNo   string
}

type InvoiceInfoResult struct {
	BizCode         string
	BizMsg          string
	OutInvoiceId    string
	OrderType       int
	InvoiceType     int
	BusinessType    int
	InvoiceTitle    string
	ApplyTime       string
	BatchInvoices   []BatchInvoice
	PayerRegisterNo string
	Remark          string
	PayerEmail      string
	ReceiverName    string
	ReceiverAddress string
	ReceiverPhone   string
}

type BatchInvoice struct {
	BatchNo          string
	SubBatchInvoices []SubBatchInvoice
	TotalCount       int
}

type SubBatchInvoice struct {
	SubBatchNo      string
	InvoiceTime     string
	Status          int
	InvoiceMaterial int
	EInvoiceUrl     string
	EInvoiceMiniUrl string
	InvoiceAmount   float64
	StartChargeSeqs string
	PickupAddress   string
}

type NotificationInvoiceChangeParam struct {
	OutInvoiceId    string
	BatchNo         string
	SubBatchNo      string
	TotalCount      int
	InvoiceMaterial int
	InvoiceTime     string
	Status          int
	EInvoiceUrl     string
	EInvoiceMiniUrl string
	StartChargeSeqs string
	InvoiceAmount   float64
	PickupAddress   string
}

type UnInvoiceSummaryParam struct {
	StartDate string
	EndDate   string
	OrderType int
}

type UnInvoiceSummaryResult struct {
	BizCode         string        `json:"BizCode"`
	BizMsg          string        `json:"BizMsg"`
	Amount          float64       `json:"Amount"`
	StartChargeSeqs []interface{} `json:"StartChargeSeqs"`
	Count           int           `json:"Count"`
	NormalAmount    float64       `json:"NormalAmount"`
	RefundAmount    float64       `json:"RefundAmount"`
	NormalCount     int           `json:"NormalCount"`
	RefundCount     int           `json:"RefundCount"`
}

//待开票订单列表

type UnInvoiceOrderParam struct {
	PageNo    int
	PageSize  int
	StartDate string
	EndDate   string
	OrderType int
}

type UnInvoiceOrderResult struct {
}

type InvoiceListParam struct {
	OutInvoiceId string
	BatchNo      string
	SubBatchNo   string
	PageNo       int
	PageSize     int
}

type InvoiceListResult struct {
	BizCode   string             `json:"BizCode"`
	PageNo    int                `json:"PageNo"`
	PageCount int                `json:"PageCount"`
	ItemSize  int                `json:"ItemSize"`
	Items     []InvoiceListItems `json:"Items"`
}

type InvoiceListItems struct {
	OutInvoiceId    string  `json:"outInvoiceId"`
	BatchNo         string  `json:"BatchNo"`
	InvoiceMaterial int     `json:"InvoiceMaterial"`
	InvoiceType     int     `json:"InvoiceType"`
	BusinessType    int     `json:"BusinessType"`
	InvoiceTitle    string  `json:"InvoiceTitle"`
	InvoiceAmount   float64 `json:"InvoiceAmount"`
	ApplyTime       string  `json:"ApplyTime"`
	Status          int     `json:"Status"`
	Content         string  `json:"Content"`
}

type InvoiceOrderParam struct {
	OutInvoiceId string
	BatchNo      string
	SubBatchNo   string
	OrderType    int
	PageNo       int
	PageSize     int
}

type InvoiceOrderResult struct {
	OutInvoiceId string `json:"OutInvoiceId"`
	BatchNo      string `json:"BatchNo"`
	SubBatchNo   string `json:"SubBatchNo"`
	OrderType    int    `json:"OrderType"`
	PageNo       int    `json:"PageNo"`
	PageSize     int    `json:"PageSize"`
}
