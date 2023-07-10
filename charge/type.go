package charge

type QueryRequest struct {
	Sig        string `json:"Sig"`
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
}

type QueryResponse struct {
	Ret  int         `json:"Ret"`
	Msg  interface{} `json:"Msg"`
	Data []byte      `json:"Data"`
	Sig  string      `json:"Sig"`
}

type ChargeResponse struct {
	Ret  int         `json:"Ret"`
	Msg  interface{} `json:"Msg"`
	Data string      `json:"Data"`
	Sig  string      `json:"Sig"`
}

type SendStarChargeParam struct {
	Data     []byte `json:"data"`
	QueryUrl string `json:"queryUrl"`
}

type SendChargeParam struct {
	Data     []byte `json:"data"`
	QueryUrl string `json:"queryUrl"`
}

type NotificationParam struct {
	Sig        string `json:"Sig"`
	Data       string `json:"Data"`
	OperatorID string `json:"OperatorID"`
	TimeStamp  string `json:"TimeStamp"`
	Seq        string `json:"Seq"`
}

type QueryEquipAuthParam struct {
	EquipAuthSeq string
	ConnectorID  string
}

type QueryEquipBusinessPolicyParam struct {
	EquipBizSeq string
	ConnectorID string
}

type PolicyInfos struct {
	StartTime   string  `json:"StartTime"`
	ElecPrice   float64 `json:"ElecPrice"`
	SevicePrice float64 `json:"SevicePrice"`
	PvType      int     `json:"PvType"`
}
type QueryEquipBusinessPolicyResult struct {
	EquipBizSeq string        `json:"EquipBizSeq"`
	ConnectorID string        `json:"ConnectorID"`
	SuccStat    int           `json:"SuccStat"`
	FailReason  int           `json:"FailReason"`
	SumPeriod   int           `json:"SumPeriod"`
	PolicyInfos []PolicyInfos `json:"PolicyInfos"`
}

type QueryEquipAuthResult struct {
	ConnectorID   string `json:"ConnectorID"`
	SuccStat      int    `json:"SuccStat"`
	FailReason    int    `json:"FailReason"`
	FailReasonMsg string `json:"FailReasonMsg"`
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
	StartChargeSeq      string          `json:"StartChargeSeq"`
	StartChargeSeqStat  int             `json:"StartChargeSeqStat"`
	ConnectorID         string          `json:"ConnectorID"`
	ConnectorStatus     int             `json:"ConnectorStatus"`
	CurrentA            float64         `json:"CurrentA"`
	CurrentB            float64         `json:"CurrentB"`
	CurrentC            float64         `json:"CurrentC"`
	VoltageA            float64         `json:"VoltageA"`
	VoltageB            float64         `json:"VoltageB"`
	VoltageC            float64         `json:"VoltageC"`
	Soc                 float64         `json:"Soc"`
	StartTime           string          `json:"StartTime"`
	EndTime             string          `json:"EndTime"`
	TotalPower          float64         `json:"TotalPower"`
	ElecMoney           float64         `json:"ElecMoney"`
	SeviceMoney         float64         `json:"SeviceMoney"`
	TotalMoney          float64         `json:"TotalMoney"`
	SumPeriod           int             `json:"SumPeriod"`
	ChargeDetails       []ChargeDetails `json:"ChargeDetails"`
	OriginalElecMoney   float64         `json:"OriginalElecMoney"`
	OriginalSeviceMoney float64         `json:"OriginalSeviceMoney"`
	OriginalMoney       float64         `json:"OriginalMoney"`
}
type QueryStopChargeParam struct {
	StartChargeSeq string
	ConnectorID    string
}
type QueryStopChargeResult struct {
	StartChargeSeq     string `json:"StartChargeSeq"`
	StartChargeSeqStat int    `json:"StartChargeSeqStat"`
	SuccStat           int    `json:"SuccStat"`
	FailReason         int    `json:"FailReason"`
}

type QueryStationsInfoParam struct {
	LastQueryTime string `json:"LastQueryTime"`
	PageNo        int64  `json:"PageNo"`
	PageSize      int64  `json:"PageSize"`
}

type QueryStationsInfoResult struct {
	PageNo       int           `json:"PageNo"`
	PageCount    int           `json:"PageCount"`
	ItemSize     int           `json:"ItemSize"`
	StationInfos []StationInfo `json:"StationInfos"`
}

type StationFeeDetail struct {
	ServiceFee     float64 `json:"ServiceFee"`
	EndTime        string  `json:"EndTime"`
	ElectricityFee float64 `json:"ElectricityFee"`
	StartTime      string  `json:"StartTime"`
	PvType         int     `json:"PvType"`
}

type OriginalStationFeeDetail struct {
	ServiceFee     float64 `json:"ServiceFee"`
	EndTime        string  `json:"EndTime"`
	ElectricityFee float64 `json:"ElectricityFee"`
	StartTime      string  `json:"StartTime"`
	PvType         int     `json:"PvType"`
}

type TariffInfo struct {
	Elements []struct {
		PriceComponents []struct {
			PriceName string `json:"PriceName"`
			PriceDesc string `json:"PriceDesc"`
		} `json:"PriceComponents"`
		Restrictions struct {
			StartTime string `json:"StartTime"`
			EndTime   string `json:"EndTime"`
			TimeType  int    `json:"TimeType"`
		} `json:"Restrictions"`
	} `json:"Elements"`
}

type BizExtParams struct {
	StubGroupId         string  `json:"StubGroupId"`
	ServiceTax          float64 `json:"ServiceTax"`
	IsBoutique          string  `json:"IsBoutique"`
	ParkingTicketPrint  int     `json:"ParkingTicketPrint"`
	RealOperatorID      string  `json:"RealOperatorID"`
	ElectricityTax      float64 `json:"ElectricityTax"`
	ParkingType         int     `json:"ParkingType"`
	EquipmentOperatorID string  `json:"EquipmentOperatorID"`
	ParkingBarrierFlag  int     `json:"ParkingBarrierFlag"`
}

type SupportingFacilityInfo struct {
	HasToilet           int `json:"HasToilet"`
	HasFood             int `json:"HasFood"`
	HasVendingMachine   int `json:"HasVendingMachine"`
	HasCanopy           int `json:"HasCanopy"`
	HasLounge           int `json:"HasLounge"`
	FreeTeaRoom         int `json:"FreeTeaRoom"`
	FreeCarWash         int `json:"FreeCarWash"`
	HasMonitor          int `json:"HasMonitor"`
	Staffing            int `json:"Staffing"`
	HasConvenienceStore int `json:"HasConvenienceStore"`
	HasCoffee           int `json:"HasCoffee"`
	HasEVService        int `json:"HasEVService"`
	HasWifi             int `json:"HasWifi"`
	HasMassageChair     int `json:"HasMassageChair"`
	HasAirConditioning  int `json:"HasAirConditioning"`
	PowerExchange       int `json:"PowerExchange"`
}

type StationInfo struct {
	StationID                string                     `json:"StationID"`
	OperatorID               string                     `json:"OperatorID"`
	EquipmentOwnerID         string                     `json:"EquipmentOwnerID"`
	StationName              string                     `json:"StationName"`
	CountryCode              string                     `json:"CountryCode"`
	AreaCode                 string                     `json:"AreaCode"`
	Address                  string                     `json:"Address"`
	StationTel               string                     `json:"StationTel"`
	ServiceTel               string                     `json:"ServiceTel"`
	StationType              int                        `json:"StationType"`
	StationStatus            int                        `json:"StationStatus"`
	ParkNums                 int                        `json:"ParkNums"`
	StationLng               float64                    `json:"StationLng"`
	StationLat               float64                    `json:"StationLat"`
	SiteGuide                string                     `json:"SiteGuide"`
	Construction             int                        `json:"Construction"`
	Pictures                 []string                   `json:"Pictures"`
	MatchCars                string                     `json:"MatchCars"`
	ParkInfo                 string                     `json:"ParkInfo"`
	BusineHours              string                     `json:"BusineHours"`
	ElectricityFee           string                     `json:"ElectricityFee"`
	ServiceFee               string                     `json:"ServiceFee"`
	ParkFee                  string                     `json:"ParkFee"`
	Payment                  string                     `json:"Payment"`
	SupportOrder             int                        `json:"SupportOrder"`
	Remark                   string                     `json:"Remark"`
	ParkingDiscountType      int                        `json:"ParkingDiscountType"`
	ParkFeeStatus            int                        `json:"ParkFeeStatus"`
	BusinessStationFeeDetail []BusinessStationFeeDetail `json:"BusinessStationFeeDetail"`
	StationFeeDetail         []StationFeeDetail         `json:"StationFeeDetail"`
	OriginalStationFeeDetail []OriginalStationFeeDetail `json:"OriginalStationFeeDetail"`
	IsEnable                 int                        `json:"IsEnable"`
	PrinterFlag              int                        `json:"PrinterFlag"`
	BarrierFlag              int                        `json:"BarrierFlag"`
	TariffInfo               TariffInfo                 `json:"TariffInfo"`
	FloorLevel               string                     `json:"FloorLevel"`
	GuideMap                 string                     `json:"GuideMap"`
	RoadInfo                 string                     `json:"RoadInfo"`
	AdminName                string                     `json:"AdminName"`
	AdminTel                 string                     `json:"AdminTel"`
	OperationWay             string                     `json:"OperationWay"`
	EnableRoaming            int                        `json:"EnableRoaming"`
	GreenEnergyFlag          int                        `json:"GreenEnergyFlag"`
	Flags                    string                     `json:"Flags"`
	EquipmentOperatorID      string                     `json:"EquipmentOperatorID"`
	OnlineTime               string                     `json:"OnlineTime"`
	StationGrade             int                        `json:"StationGrade"`
	BizExtParams             BizExtParams               `json:"BizExtParams"`
	SupportingFacilityInfo   SupportingFacilityInfo     `json:"SupportingFacilityInfo"`
	PlaceHolderType          int                        `json:"PlaceHolderType"`
	EquipmentInfos           []EquipmentInfos           `json:"EquipmentInfos"`
}

type BusinessStationFeeDetail struct {
	ServiceFee     float64 `json:"ServiceFee"`
	EndTime        string  `json:"EndTime"`
	ElectricityFee float64 `json:"ElectricityFee"`
	StartTime      string  `json:"StartTime"`
	PvType         int     `json:"PvType"`
}
type EquipmentInfos struct {
	EquipmentID      string  `json:"EquipmentID"`
	ManufacturerName string  `json:"ManufacturerName"`
	EquipmentModel   string  `json:"EquipmentModel"`
	ProductionDate   string  `json:"ProductionDate"`
	EquipmentType    int     `json:"EquipmentType"`
	EquipmentLng     float64 `json:"EquipmentLng"`
	EquipmentLat     float64 `json:"EquipmentLat"`
	Power            float64 `json:"Power"`
	EquipmentName    string  `json:"EquipmentName"`
	PowerType        int     `json:"PowerType"`
	BizExtParams     struct {
		ChargePointId string `json:"ChargePointId"`
	} `json:"BizExtParams"`
	ConnectorInfos []ConnectorInfo `json:"ConnectorInfos"`
	ManufacturerID string          `json:"ManufacturerID,omitempty"`
}

type ConnectorInfo struct {
	ConnectorID        string  `json:"ConnectorID"`
	ConnectorName      string  `json:"ConnectorName"`
	ConnectorType      int     `json:"ConnectorType"`
	VoltageUpperLimits int     `json:"VoltageUpperLimits"`
	VoltageLowerLimits int     `json:"VoltageLowerLimits"`
	Current            int     `json:"Current"`
	Power              float64 `json:"Power"`
	ParkNo             string  `json:"ParkNo"`
	NationalStandard   int     `json:"NationalStandard"`
	ParkingLockFlag    int     `json:"ParkingLockFlag"`
	Reservable         int     `json:"Reservable"`
	BluetoothKey       string  `json:"BluetoothKey"`
	VinSupport         int     `json:"VinSupport"`
	BizExtParams       struct {
		ModelNo    string `json:"ModelNo"`
		StubStatus string `json:"StubStatus"`
		StubId     string `json:"StubId"`
		UpdateType string `json:"UpdateType"`
		FrameWork  string `json:"FrameWork"`
	} `json:"BizExtParams"`
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

type QueryTokenResult struct {
	OperatorID         string `json:"OperatorID"`
	SuccStat           int    `json:"SuccStat"`
	AccessToken        string `json:"AccessToken"`
	TokenAvailableTime int    `json:"TokenAvailableTime"`
	FailReason         int    `json:"FailReason"`
}

type StationStatusInfos struct {
	StationID            string                `json:"StationID"`
	ConnectorStatusInfos []ConnectorStatusInfo `json:"ConnectorStatusInfos"`
}

type ConnectorStatusInfo struct {
	ConnectorID string `json:"ConnectorID"`
	Status      int    `json:"Status"`
	ParkStatus  int    `json:"ParkStatus"`
	LockStatus  int    `json:"LockStatus"`
	StatusCode  string `json:"StatusCode,omitempty"`
}

type QueryStationStatusResult struct {
	Total              int                  `json:"Total"`
	StationStatusInfos []StationStatusInfos `json:"StationStatusInfos"`
}

type QueryTokenParam struct {
	OperatorSecret string `json:"OperatorSecret"`
	OperatorID     string `json:"OperatorID"`
}

type QueryTokenReq struct {
	OperatorID     string `json:"OperatorId" form:"OperatorId" binding:"required"`
	OperatorSecret string `json:"OperatorSecret" form:"OperatorSecret" binding:"required"`
}

type QueryTokenResp struct {
	OperatorID         string `json:"OperatorId"`
	SuccStat           int64  `json:"SuccStat"`
	AccessToken        string `json:"AccessToken"`
	TokenAvailableTime int64  `json:"TokenAvailableTime"`
	FailReason         int64  `json:"FailReason"`
}

type NotificationStationStatusParam struct {
	StationStatusInfo ConnectorStatusInfo `json:"ConnectorStatusInfo"`
}

type NotificationStationInfoParam struct {
	StationInfo StationInfo `json:"StationInfo"`
}

type NotificationStationStatusResult struct {
	Status string `json:"Status"`
}

type NotificationStartChargeResultParam struct {
	StartChargeSeq     string
	StartChargeSeqStat int64
	ConnectorID        string
	StartTime          string
	IdentCode          string
}

type ChargeDetails struct {
	DetailPower       float64 `json:"DetailPower"`
	ElecPrice         float64 `json:"ElecPrice"`
	SevicePrice       float64 `json:"SevicePrice"`
	DetailElecMoney   float64 `json:"DetailElecMoney"`
	DetailSeviceMoney float64 `json:"DetailSeviceMoney"`
	DetailStartTime   string  `json:"DetailStartTime"`
	DetailEndTime     string  `json:"DetailEndTime"`
	DetailType        int     `json:"DetailType"`
}

type NotificationEquipChargeStatusParamV2 struct {
	StartChargeSeq     string          `json:"StartChargeSeq"`
	StartChargeSeqStat int             `json:"StartChargeSeqStat"`
	ConnectorID        string          `json:"ConnectorID"`
	ConnectorStatus    int             `json:"ConnectorStatus"`
	CurrentA           float64         `json:"CurrentA"`
	CurrentB           float64         `json:"CurrentB"`
	CurrentC           float64         `json:"CurrentC"`
	VoltageA           float64         `json:"VoltageA"`
	VoltageB           float64         `json:"VoltageB"`
	VoltageC           float64         `json:"VoltageC"`
	Soc                float64         `json:"Soc"`
	StartTime          string          `json:"StartTime"`
	EndTime            string          `json:"EndTime"`
	TotalPower         float64         `json:"TotalPower"`
	ElecMoney          float64         `json:"ElecMoney"`
	SeviceMoney        float64         `json:"SeviceMoney"`
	TotalMoney         float64         `json:"TotalMoney"`
	SumPeriod          int             `json:"SumPeriod"`
	ChargeDetails      []ChargeDetails `json:"ChargeDetails"`
}

type TariffChargeDetails struct {
	TariffChargeName string `json:"TariffChargeName"`
	TariffChargeDesc string `json:"TariffChargeDesc"`
}

type NotificationEquipChargeStatusParam struct {
	StartChargeSeq        string                        `json:"StartChargeSeq"`
	StartChargeSeqStat    int                           `json:"StartChargeSeqStat"`
	ConnectorID           string                        `json:"ConnectorID"`
	ConnectorStatus       int                           `json:"ConnectorStatus"`
	CurrentA              float64                       `json:"CurrentA"`
	CurrentB              float64                       `json:"CurrentB"`
	CurrentC              float64                       `json:"CurrentC"`
	VoltageA              float64                       `json:"VoltageA"`
	VoltageB              float64                       `json:"VoltageB"`
	VoltageC              float64                       `json:"VoltageC"`
	Soc                   float64                       `json:"Soc"`
	StartTime             string                        `json:"StartTime"`
	EndTime               string                        `json:"EndTime"`
	TotalPower            float64                       `json:"TotalPower"`
	ElecMoney             float64                       `json:"ElecMoney"`
	SeviceMoney           float64                       `json:"SeviceMoney"`
	TotalMoney            float64                       `json:"TotalMoney"`
	SumPeriod             int                           `json:"SumPeriod"`
	ChargeDetails         []ChargeDetails               `json:"ChargeDetails"`
	PWM                   int                           `json:"PWM"`
	CurrentFrequency      int                           `json:"CurrentFrequency"`
	ExpectEndTime         string                        `json:"ExpectEndTime"`
	TariffChargeInfo      TariffChargeInfo              `json:"TariffChargeInfo"`
	BizExtParams          EquipChargeStatusBizExtParams `json:"BizExtParams"`
	OrderID               string                        `json:"OrderID"`
	StartType             int                           `json:"StartType"`
	LeftTime              int                           `json:"LeftTime"`
	VoltageCar            int                           `json:"VoltageCar"`
	CurrentCar            int                           `json:"CurrentCar"`
	OriginalElecMoney     float64                       `json:"OriginalElecMoney"`
	OriginalSeviceMoney   float64                       `json:"OriginalSeviceMoney"`
	OriginalMoney         float64                       `json:"OriginalMoney"`
	OriginalChargeDetails []OriginalChargeDetails       `json:"OriginalChargeDetails"`
}

type TariffChargeInfo struct {
	TariffChargeDetails []TariffChargeDetails `json:"TariffChargeDetails"`
}

type EquipChargeStatusBizExtParams struct {
	CTS                string  `json:"CTS"`
	ChargeOperatorID   string  `json:"ChargeOperatorID"`
	TempStub           int     `json:"TempStub"`
	OrderID            string  `json:"OrderID"`
	StationID          string  `json:"StationID"`
	CurrentCar         int     `json:"CurrentCar"`
	LeftTime           int     `json:"LeftTime"`
	TargetOperatorID   string  `json:"TargetOperatorID"`
	StartType          int     `json:"StartType"`
	StationType        int     `json:"StationType"`
	EquipmentID        string  `json:"EquipmentID"`
	ElectricMeterEnd   float64 `json:"ElectricMeterEnd"`
	TotalPower         float64 `json:"TotalPower"`
	TempGun            int     `json:"TempGun"`
	AreaCode           string  `json:"AreaCode"`
	ElectricMeterStart int     `json:"ElectricMeterStart"`
	VoltageCar         int     `json:"VoltageCar"`
	ConnectorID        string  `json:"ConnectorID"`
}

type OriginalChargeDetails struct {
	DetailPower       float64 `json:"DetailPower"`
	ElecPrice         float64 `json:"ElecPrice"`
	SevicePrice       float64 `json:"SevicePrice"`
	DetailElecMoney   float64 `json:"DetailElecMoney"`
	DetailSeviceMoney float64 `json:"DetailSeviceMoney"`
	DetailStartTime   string  `json:"DetailStartTime"`
	DetailEndTime     string  `json:"DetailEndTime"`
}

type NotificationStopChargeResultParam struct {
	StartChargeSeq     string
	StartChargeSeqStat int64
	ConnectorID        string
	SuccStat           int64
	FailReason         int64
}

type NotificationChargeOrderInfoParam struct {
	StartChargeSeq           string                      `json:"StartChargeSeq"`
	ConnectorId              string                      `json:"ConnectorID"`
	StartTime                string                      `json:"StartTime"`
	EndTime                  string                      `json:"EndTime"`
	TotalPower               float64                     `json:"TotalPower"`
	TotalElecMoney           float64                     `json:"TotalElecMoney"`
	TotalSeviceMoney         float64                     `json:"TotalSeviceMoney"`
	TotalMoney               float64                     `json:"TotalMoney"`
	StopReason               int                         `json:"StopReason"`
	SumPeriod                int                         `json:"SumPeriod"`
	ChargeDetails            []ChargeDetails             `json:"ChargeDetails"`
	TariffChargeInfo         TariffChargeInfo            `json:"TariffChargeInfo"`
	BizExtParams             ChargeOrderInfoBizExtParams `json:"BizExtParams"`
	OrderID                  string                      `json:"OrderID"`
	StationOwnerType         int                         `json:"StationOwnerType"`
	StartType                int                         `json:"StartType"`
	TotalBusinessElecMoney   float64                     `json:"TotalBusinessElecMoney"`
	TotalBusinessSeviceMoney float64                     `json:"TotalBusinessSeviceMoney"`
	TotalBusinessMoney       float64                     `json:"TotalBusinessMoney"`
	BusinessChargeDetails    []BusinessChargeDetails     `json:"BusinessChargeDetails"`
	TotalOriginalElecMoney   float64                     `json:"TotalOriginalElecMoney"`
	TotalOriginalSeviceMoney float64                     `json:"TotalOriginalSeviceMoney"`
	TotalOriginalMoney       float64                     `json:"TotalOriginalMoney"`
	OriginalChargeDetails    []OriginalChargeDetails     `json:"OriginalChargeDetails"`
	EndSOC                   int                         `json:"EndSOC"`
	ElectricStart            int                         `json:"ElectricStart"`
	ElectricEnd              float64                     `json:"ElectricEnd"`
	UserID                   string                      `json:"UserID"`
}

type ChargeOrderInfoBizExtParams struct {
	CTS                string  `json:"CTS"`
	ChargeOperatorID   string  `json:"ChargeOperatorID"`
	OrderID            string  `json:"OrderID"`
	TargetOperatorID   string  `json:"TargetOperatorID"`
	StartType          int     `json:"StartType"`
	StationType        int     `json:"StationType"`
	Soc                int     `json:"Soc"`
	IsMergeOrder       int     `json:"IsMergeOrder"`
	ElectricMeterEnd   float64 `json:"ElectricMeterEnd"`
	TotalPower         float64 `json:"TotalPower"`
	TotalFeeInfo       string  `json:"TotalFeeInfo"`
	EndType            int     `json:"EndType"`
	CTL                string  `json:"CTL"`
	ElectricMeterStart int     `json:"ElectricMeterStart"`
	CTP                int     `json:"CTP"`
}

type BusinessChargeDetails struct {
	DetailPower       float64 `json:"DetailPower"`
	ElecPrice         float64 `json:"ElecPrice"`
	SevicePrice       float64 `json:"SevicePrice"`
	DetailElecMoney   float64 `json:"DetailElecMoney"`
	DetailSeviceMoney float64 `json:"DetailSeviceMoney"`
	DetailStartTime   string  `json:"DetailStartTime"`
	DetailEndTime     string  `json:"DetailEndTime"`
	DetailType        int     `json:"DetailType"`
}
