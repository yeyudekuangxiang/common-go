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
	StartChargeSeq     string  `json:"StartChargeSeq"`
	StartChargeSeqStat int     `json:"StartChargeSeqStat"`
	ConnectorID        string  `json:"ConnectorID"`
	ConnectorStatus    int     `json:"ConnectorStatus"`
	CurrentA           float64 `json:"CurrentA"`
	CurrentB           float64 `json:"CurrentB"`
	CurrentC           float64 `json:"CurrentC"`
	VoltageA           float64 `json:"VoltageA"`
	VoltageB           float64 `json:"VoltageB"`
	VoltageC           float64 `json:"VoltageC"`
	Soc                float64 `json:"Soc"`
	StartTime          string  `json:"StartTime"`
	EndTime            string  `json:"EndTime"`
	TotalPower         float64 `json:"TotalPower"`
	ElecMoney          float64 `json:"ElecMoney"`
	SeviceMoney        float64 `json:"SeviceMoney"`
	TotalMoney         float64 `json:"TotalMoney"`
	SumPeriod          int     `json:"SumPeriod"`
	ChargeDetails      []struct {
		DetailPower       float64 `json:"DetailPower"`
		ElecPrice         float64 `json:"ElecPrice"`
		SevicePrice       float64 `json:"SevicePrice"`
		DetailElecMoney   float64 `json:"DetailElecMoney"`
		DetailSeviceMoney float64 `json:"DetailSeviceMoney"`
		DetailStartTime   string  `json:"DetailStartTime"`
		DetailEndTime     string  `json:"DetailEndTime"`
		DetailType        int     `json:"DetailType"`
	} `json:"ChargeDetails"`
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
	PageNo       int `json:"PageNo"`
	PageCount    int `json:"PageCount"`
	ItemSize     int `json:"ItemSize"`
	StationInfos []struct {
		StationID                string   `json:"StationID"`
		OperatorID               string   `json:"OperatorID"`
		EquipmentOwnerID         string   `json:"EquipmentOwnerID"`
		StationName              string   `json:"StationName"`
		CountryCode              string   `json:"CountryCode"`
		AreaCode                 string   `json:"AreaCode"`
		Address                  string   `json:"Address"`
		StationTel               string   `json:"StationTel"`
		ServiceTel               string   `json:"ServiceTel"`
		StationType              int      `json:"StationType"`
		StationStatus            int      `json:"StationStatus"`
		ParkNums                 int      `json:"ParkNums"`
		StationLng               float64  `json:"StationLng"`
		StationLat               float64  `json:"StationLat"`
		SiteGuide                string   `json:"SiteGuide"`
		Construction             int      `json:"Construction"`
		Pictures                 []string `json:"Pictures"`
		MatchCars                string   `json:"MatchCars"`
		ParkInfo                 string   `json:"ParkInfo"`
		BusineHours              string   `json:"BusineHours"`
		ElectricityFee           string   `json:"ElectricityFee"`
		ServiceFee               string   `json:"ServiceFee"`
		ParkFee                  string   `json:"ParkFee"`
		Payment                  string   `json:"Payment"`
		SupportOrder             int      `json:"SupportOrder"`
		Remark                   string   `json:"Remark"`
		ParkingDiscountType      int      `json:"ParkingDiscountType"`
		ParkFeeStatus            int      `json:"ParkFeeStatus"`
		BusinessStationFeeDetail []struct {
			ServiceFee     float64 `json:"ServiceFee"`
			EndTime        string  `json:"EndTime"`
			ElectricityFee float64 `json:"ElectricityFee"`
			StartTime      string  `json:"StartTime"`
			PvType         int     `json:"PvType"`
		} `json:"BusinessStationFeeDetail"`
		StationFeeDetail []struct {
			ServiceFee     float64 `json:"ServiceFee"`
			EndTime        string  `json:"EndTime"`
			ElectricityFee float64 `json:"ElectricityFee"`
			StartTime      string  `json:"StartTime"`
			PvType         int     `json:"PvType"`
		} `json:"StationFeeDetail"`
		OriginalStationFeeDetail []struct {
			ServiceFee     float64 `json:"ServiceFee"`
			EndTime        string  `json:"EndTime"`
			ElectricityFee float64 `json:"ElectricityFee"`
			StartTime      string  `json:"StartTime"`
			PvType         int     `json:"PvType"`
		} `json:"OriginalStationFeeDetail"`
		IsEnable    int `json:"IsEnable"`
		PrinterFlag int `json:"PrinterFlag"`
		BarrierFlag int `json:"BarrierFlag"`
		TariffInfo  struct {
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
		} `json:"TariffInfo"`
		FloorLevel          string `json:"FloorLevel"`
		GuideMap            string `json:"GuideMap"`
		RoadInfo            string `json:"RoadInfo"`
		AdminName           string `json:"AdminName"`
		AdminTel            string `json:"AdminTel"`
		OperationWay        string `json:"OperationWay"`
		EnableRoaming       int    `json:"EnableRoaming"`
		GreenEnergyFlag     int    `json:"GreenEnergyFlag"`
		Flags               string `json:"Flags"`
		EquipmentOperatorID string `json:"EquipmentOperatorID"`
		OnlineTime          string `json:"OnlineTime"`
		StationGrade        int    `json:"StationGrade"`
		BizExtParams        struct {
			StubGroupId         string  `json:"StubGroupId"`
			ServiceTax          float64 `json:"ServiceTax"`
			IsBoutique          string  `json:"IsBoutique"`
			ParkingTicketPrint  int     `json:"ParkingTicketPrint"`
			RealOperatorID      string  `json:"RealOperatorID"`
			ElectricityTax      float64 `json:"ElectricityTax"`
			ParkingType         int     `json:"ParkingType"`
			EquipmentOperatorID string  `json:"EquipmentOperatorID"`
			ParkingBarrierFlag  int     `json:"ParkingBarrierFlag"`
		} `json:"BizExtParams"`
		SupportingFacilityInfo struct {
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
		} `json:"SupportingFacilityInfo"`
		PlaceHolderType int `json:"PlaceHolderType"`
		EquipmentInfos  []struct {
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
			ConnectorInfos []struct {
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
			} `json:"ConnectorInfos"`
			ManufacturerID string `json:"ManufacturerID,omitempty"`
		} `json:"EquipmentInfos"`
	} `json:"StationInfos"`
}

type QueryStationsInfoResultV2 struct {
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

type QueryTokenResult struct {
	OperatorID         string `json:"OperatorID"`
	SuccStat           int    `json:"SuccStat"`
	AccessToken        string `json:"AccessToken"`
	TokenAvailableTime int    `json:"TokenAvailableTime"`
	FailReason         int    `json:"FailReason"`
}

type QueryStationStatusResult struct {
	Total              int `json:"Total"`
	StationStatusInfos []struct {
		StationID            string `json:"StationID"`
		ConnectorStatusInfos []struct {
			ConnectorID string `json:"ConnectorID"`
			Status      int    `json:"Status"`
			ParkStatus  int    `json:"ParkStatus"`
			LockStatus  int    `json:"LockStatus"`
			StatusCode  string `json:"StatusCode,omitempty"`
		} `json:"ConnectorStatusInfos"`
	} `json:"StationStatusInfos"`
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
	StationStatusInfo StationStatusInfo
}

type StationStatusInfo struct {
	ConnectorID string `json:"ConnectorID"`
	Status      int64  `json:"Status"`
	ParkStatus  int64  `json:"ParkStatus"`
	LockStatus  int64  `json:"LockStatus"`
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

type NotificationEquipChargeStatusParam struct {
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

type NotificationStopChargeResultParam struct {
	StartChargeSeq     string
	StartChargeSeqStat int64
	ConnectorID        string
	SuccStat           int64
	FailReason         int64
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
