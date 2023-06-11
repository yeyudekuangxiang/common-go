package starcharge

type SendStarChargeParam struct {
	Data     []byte `json:"data"`
	QueryUrl string `json:"queryUrl"`
}
type starChargeResponse struct {
	Ret  int    `json:"Ret"`
	Msg  string `json:"Msg"`
	Data []byte `json:"Data"`
	Sig  string `json:"Sig"`
}

type QueryEquipAuthParam struct {
	EquipBizSeq string
	ConnectorID string
}

type QueryEquipAuthResult struct {
	EquipAuthSeq string `json:"EquipAuthSeq"`
	ConnectorID  string `json:"ConnectorID"`
	SuccStat     int    `json:"SuccStat"`
	FailReason   int    `json:"FailReason"`
}

type QueryStartChargeParam struct {
	StartChargeSeq string //订单号
	ConnectorID    string
	QRCode         string //二维码其他信息
}

type QueryStartChargeResult struct {
	StartChargeSeq     string `json:"StartChargeSeq"`
	StartChargeSeqStat int    `json:"StartChargeSeqStat"` //充电订单状态
	ConnectorID        string `json:"ConnectorID"`
	SuccStat           int    `json:"SuccStat"`
	FailReason         int    `json:"FailReason"`
}

type QueryEquipChargeStatusParam struct {
	StartChargeSeq string `json:"StartChargeSeq"`
}
type QueryEquipChargeStatusResult struct {
	StartChargeSeq     string  `json:"StartChargeSeq"`
	StartChargeSeqStat int64   `json:"StartChargeSeqStat"`
	ConnectorID        string  `json:"ConnectorID"`
	StartTime          string  `json:"StartTime"`
	IdentCode          string  `json:"IdentCode"`
	ConnectorStatus    int64   `json:"ConnectorStatus"`
	CurrentA           float64 `json:"CurrentA"`
	CurrentB           float64 `json:"CurrentB"`
	CurrentC           float64 `json:"CurrentC"`
	VoltageA           float64 `json:"VoltageA"`
	VoltageB           float64 `json:"VoltageB"`
	VoltageC           float64 `json:"VoltageC"`
	Soc                float64 `json:"Soc"`
	EndTime            string  `json:"EndTime"`
	TotalPower         float64 `json:"TotalPower"`
	ElecMoney          float64 `json:"ElecMoney"`
	ServiceMoney       float64 `json:"ServiceMoney"`
	TotalMoney         float64 `json:"TotalMoney"`
	StopReason         int64   `json:"StopReason"`
	SumPeriod          int64   `json:"SumPeriod"`
	ChargeDetails      string  `json:"ChargeDetails"`
}

type QueryStopChargeParam struct {
	StartChargeSeq string
	ConnectorID    string
}
type QueryStopChargeResult struct {
	StartChargeSeq     string `json:"StartChargeSeq"`
	StartChargeSeqStat int64  `json:"StartChargeSeqStat"`
	SuccStat           int    `json:"SuccStat"`
	FailReason         int    `json:"FailReason"`
}
