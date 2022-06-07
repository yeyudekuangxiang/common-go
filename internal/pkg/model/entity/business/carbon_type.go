package business

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"time"
)

type CarbonType string

const (
	CarbonTypeOnlineMeeting        CarbonType = "OnlineMeeting"        //线上会议
	CarbonTypeSaveWaterElectricity CarbonType = "SaveWaterElectricity" //节水节电
	CarbonTypePublicTransport      CarbonType = "PublicTransport"      //公交地铁
	CarbonTypeEvCar                CarbonType = "EvCar"                //电动车 电车充电
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

// CalcOnlineMeeting 根据会议时长计算获得多少碳积分
func (t CarbonType) CalcOnlineMeeting(m time.Duration) decimal.Decimal {
	panic("请配置我")
}

// CalcSaveWaterElectricity 根绝节水和节电量计算获得多少碳积分 水的单位升 电的单位度
func (t CarbonType) CalcSaveWaterElectricity(water int64, electricity int64) decimal.Decimal {
	panic("请配置我")
}

// CalcPublicTransport 根绝公交和地铁的距离计算获得多少碳积分 单位都是公里
func (t CarbonType) CalcPublicTransport(bus int64, metro int64) decimal.Decimal {
	panic("请配置我")
}

// CalcEvCar 根绝电车充电量计算获得多少碳积分 单位度
func (t CarbonType) CalcEvCar(electricity int64) decimal.Decimal {
	panic("请配置我")
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
	MeetingDuration time.Duration `json:"meetingDuration"` //会议时长
	StartTime       model.Time    `json:"startTime"`
	EndTime         model.Time    `json:"endTime"`
}

func (c CarbonTypeInfoOnlineMeeting) JSON() string {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return string(data)
}

type CarbonTypeInfoSaveWaterElectricity struct {
	Water       int64 `json:"water"`       //数量 升
	Electricity int64 `json:"electricity"` //电量 度
}

func (c CarbonTypeInfoSaveWaterElectricity) JSON() string {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return string(data)
}

type CarbonTypeInfoPublicTransport struct {
	Bus   int64 //公里
	Metro int64 //公里
}

func (c CarbonTypeInfoPublicTransport) JSON() string {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return string(data)
}

type CarbonTypeInfoEvCar struct {
	Electricity int64 //度
}

func (c CarbonTypeInfoEvCar) JSON() string {
	data, err := json.Marshal(c)
	if err != nil {
		app.Logger.Error(err)
	}
	return string(data)
}
