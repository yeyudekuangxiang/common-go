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
