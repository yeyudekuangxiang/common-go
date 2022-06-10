package business

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

// PointSetting 碳积分和积分汇率配置json字符串
type PointSetting string

type PointSettingRate interface {
	PointSetting() PointSetting
}

// SaveWaterElectricityExchangeRate 节水节电碳积分和积分兑换比例
type SaveWaterElectricityExchangeRate struct {
	Water       PointExchangeRate `json:"water"`
	Electricity PointExchangeRate `json:"electricity"`
}

func (o SaveWaterElectricityExchangeRate) PointSetting() PointSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointSetting(data)
}

// PublicTransportExchangeRate 公交地铁碳积分和积分兑换比例
type PublicTransportExchangeRate struct {
	Bus   PointExchangeRate `json:"bus"`
	Metro PointExchangeRate `json:"metro"`
}

func (o PublicTransportExchangeRate) PointSetting() PointSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointSetting(data)
}

// OnlineMeetingExchangeRate 线上会议碳积分和积分兑换比例
type OnlineMeetingExchangeRate struct {
	OneCity  PointExchangeRate `json:"oneCity"`  //异地会议
	ManyCity PointExchangeRate `json:"manyCity"` //同城会议
}

func (o OnlineMeetingExchangeRate) PointSetting() PointSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointSetting(data)
}

// PointExchangeRate 碳积分和积分兑换比例
type PointExchangeRate struct {
	CarbonCredit decimal.Decimal `json:"carbonCredit"`
	Point        int             `json:"point"`
}

func (o PointExchangeRate) PointSetting() PointSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointSetting(data)
}

// Calc 根据碳积分计算可得积分
func (o PointExchangeRate) Calc(carbonCredit decimal.Decimal) int {
	return int(carbonCredit.Div(carbonCredit).Mul(decimal.NewFromInt(int64(o.Point))).IntPart())
}
