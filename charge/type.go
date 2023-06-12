package charge

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

type NotificationParam struct {
	Sig        string `json:"Sig"`
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
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

type QueryStationsInfoParam struct {
	LastQueryTime string `json:"LastQueryTime"`
	PageNo        int64  `json:"PageNo"`
	PageSize      int64  `json:"PageSize"`
}

type QueryStationsInfoResult struct {
	ItemSize     int            `json:"ItemSize"`
	PageCount    int            `json:"pageCount"`
	PageNo       int            `json:"PageNo"`
	StationInfos []StationInfos `json:"StationInfos"`
}

type StationInfos struct {
	OperationID      string           `json:"OperationID"`
	StationID        string           `json:"StationID"`
	StationName      string           `json:"StationName"`
	EquipmentOwnerID string           `json:"EquipmentOwnerID"`
	CountryCode      string           `json:"CountryCode"`
	AreaCode         string           `json:"AreaCode"`
	Address          string           `json:"Address"`
	ServiceTel       string           `json:"ServiceTel"`
	StationType      int64            `json:"StationType"`
	StationStatus    int64            `json:"StationStatus"`
	ParkNums         int64            `json:"ParkNums"`
	StationLng       float64          `json:"StationLng"`
	StationLat       float64          `json:"StationLat"`
	Construction     int64            `json:"Construction"`
	SiteGuide        string           `json:"SiteGuide"`
	MatchCars        string           `json:"MatchCars"`
	ParkInfo         string           `json:"ParkInfo"`
	BusineHours      string           `json:"BusineHours"`
	ElectricityFee   string           `json:"ElectricityFee"`
	ServiceFee       string           `json:"ServiceFee"`
	ParkFee          string           `json:"ParkFee"`
	Payment          string           `json:"Payment"`
	SupportOrder     int64            `json:"SupportOrder"`
	Remark           string           `json:"Remark"`
	Pictures         []string         `json:"Pictures"`
	EquipmentInfos   []EquipmentInfos `json:"EquipmentInfos"`
}

type EquipmentInfos struct {
	EquipmentID    string           `json:"EquipmentID"`
	ManufacturerID string           `json:"ManufacturerID"`
	EquipmentModel string           `json:"EquipmentModel"`
	EquipmentType  int64            `json:"EquipmentType"`
	Power          string           `json:"Power"`
	EquipmentName  string           `json:"EquipmentName"`
	ConnectorInfos []ConnectorInfos `json:"ConnectorInfos"`
}

type ConnectorInfos struct {
	ConnectorID        string  `json:"ConnectorID"`
	ConnectorName      string  `json:"ConnectorName"`
	ConnectorType      int64   `json:"ConnectorType"`
	VoltageUpperLimits int64   `json:"VoltageUpperLimits"`
	VoltageLowerLimits int64   `json:"VoltageLowerLimits"`
	Current            int64   `json:"Current"`
	NationalStandard   int64   `json:"NationalStandard"`
	Power              float64 `json:"Power"`
	ParkNo             string  `json:"ParkNo"`
}

type QueryStationStatusParam struct {
	StationIDs []string `json:"StationIDs"`
}

type QueryStationStatusResult struct {
}

type NotificationChargeOrderInfoParam struct {
	StartChargeSeq   string
	ConnectorId      string
	StartTime        string
	EndTime          string
	TotalPower       float64
	TotalElecMoney   float64
	TotalSeviceMoney float64
	TotalMoney       float64
	StopReason       int64
	SumPeriod        int64
	ChargeDetails    string
}

type NotificationEquipChargeStatusParam struct {
	StartChargeSeq     string
	StartChargeSeqStat int64
	ConnectorID        string
	StartTime          string
	IdentCode          string
	ConnectorStatus    int64
	CurrentA           float64
	CurrentB           float64
	CurrentC           float64
	VoltageA           float64
	VoltageB           float64
	VoltageC           float64
	Soc                float64
	EndTime            string
	TotalPower         float64
	ElecMoney          float64
	ServiceMoney       float64
	TotalMoney         float64
	StopReason         int64
	SumPeriod          int64
	ChargeDetails      string
}
