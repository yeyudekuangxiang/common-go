package business

import (
	"encoding/json"
	"mio/internal/pkg/core/app"
	"time"
)

type CarbonType string

const (
	CarbonTypeOnlineMeeting        CarbonType = "OnlineMeeting"        //线上会议
	CarbonTypeSaveWaterElectricity CarbonType = "SaveWaterElectricity" //节水节电
	CarbonTypePublicTransport      CarbonType = "PublicTransport"      //公交地铁
	CarbonTypeEvCar                CarbonType = "EvCar"                //电动车 电车充电
	CarbonTypeOEP                  CarbonType = "OEP"                  //光盘行动
	CarbonTypeGreenBusinessTrip    CarbonType = "GreenBusinessTrip"    //低碳旅行
)

// Text 展示给用户看的
func (t CarbonType) Text() string {
	switch t {
	case CarbonTypeOnlineMeeting:
		return "线上会议"
	case CarbonTypeSaveWaterElectricity:
		return "节水节电"
	case CarbonTypePublicTransport:
		return "公交地铁"
	case CarbonTypeEvCar:
		return "电车充电"
	}
	return "未知类型"
}

// RealText 展示给管理员看的
func (t CarbonType) RealText() string {
	switch t {
	case CarbonTypeOnlineMeeting:
		return "线上会议"
	case CarbonTypeSaveWaterElectricity:
		return "节水节电"
	case CarbonTypePublicTransport:
		return "公交地铁"
	case CarbonTypeEvCar:
		return "电车充电"
	}
	return "未知类型"
}
func (t CarbonType) PointType() PointType {
	switch t {
	case CarbonTypeOnlineMeeting:
		return PointTypeOnlineMeeting
	case CarbonTypeSaveWaterElectricity:
		return PointTypeSaveWaterElectricity
	case CarbonTypePublicTransport:
		return PointTypePublicTransport
	case CarbonTypeEvCar:
		return PointTypeEvCar
	}
	return ""
}

type CarbonTypeInfo string

func (info CarbonTypeInfo) OnlineMeeting() (CarbonTypeInfoOnlineMeeting, error) {
	meeting := CarbonTypeInfoOnlineMeeting{}
	return meeting, json.Unmarshal([]byte(info), &meeting)
}
func (info CarbonTypeInfo) SaveWaterElectricity() (CarbonTypeInfoSaveWaterElectricity, error) {
	we := CarbonTypeInfoSaveWaterElectricity{}
	return we, json.Unmarshal([]byte(info), &we)
}
func (info CarbonTypeInfo) PublicTransport() (CarbonTypeInfoPublicTransport, error) {
	pt := CarbonTypeInfoPublicTransport{}
	return pt, json.Unmarshal([]byte(info), &pt)
}
func (info CarbonTypeInfo) EvCar() (CarbonTypeInfoEvCar, error) {
	ec := CarbonTypeInfoEvCar{}
	return ec, json.Unmarshal([]byte(info), &ec)
}

// CarbonTypeInfoOnlineMeeting 会议信息
type CarbonTypeInfoOnlineMeeting struct {
	OneCityDuration  time.Duration `json:"OneCityDuration"`  //同城在线会议时长
	ManyCityDuration time.Duration `json:"manyCityDuration"` //异地在线会议时长
}

func (c CarbonTypeInfoOnlineMeeting) CarbonTypeInfo() CarbonTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return CarbonTypeInfo(data)
}
func (c CarbonTypeInfoOnlineMeeting) PointTypeInfo() PointTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return PointTypeInfo(data)
}

type CarbonTypeInfoSaveWaterElectricity struct {
	Water       int64 `json:"water"`       //数量 升
	Electricity int64 `json:"electricity"` //电量 度
}

func (c CarbonTypeInfoSaveWaterElectricity) CarbonTypeInfo() CarbonTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return CarbonTypeInfo(data)
}
func (c CarbonTypeInfoSaveWaterElectricity) PointTypeInfo() PointTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return PointTypeInfo(data)
}

type CarbonTypeInfoPublicTransport struct {
	Bus   float64 //公里
	Metro float64 //公里
	Step  float64
	Bike  float64
}

func (c CarbonTypeInfoPublicTransport) CarbonTypeInfo() CarbonTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return CarbonTypeInfo(data)
}
func (c CarbonTypeInfoPublicTransport) PointTypeInfo() PointTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return PointTypeInfo(data)
}

type CarbonTypeInfoEvCar struct {
	Electricity float64 //度
}

func (c CarbonTypeInfoEvCar) CarbonTypeInfo() CarbonTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return CarbonTypeInfo(data)
}
func (c CarbonTypeInfoEvCar) PointTypeInfo() PointTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return PointTypeInfo(data)
}

type CarbonTypeInfoOEP struct {
	Voucher string
}

func (c CarbonTypeInfoOEP) CarbonTypeInfo() CarbonTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return CarbonTypeInfo(data)
}
func (c CarbonTypeInfoOEP) PointTypeInfo() PointTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return PointTypeInfo(data)
}

type CarbonTypeInfoGreenBusinessTrip struct {
	Distance float64
	From     string
	To       string
}

func (c CarbonTypeInfoGreenBusinessTrip) CarbonTypeInfo() CarbonTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return CarbonTypeInfo(data)
}
func (c CarbonTypeInfoGreenBusinessTrip) PointTypeInfo() PointTypeInfo {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return PointTypeInfo(data)
}
