package business

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

// PointRateSetting 碳积分和积分汇率配置json字符串
type PointRateSetting string

type IPointRateSetting interface {
	PointRateSetting() PointRateSetting
}

// PointRateSaveWaterElectricity 节水节电碳积分和积分兑换比例
type PointRateSaveWaterElectricity struct {
	Water       PointRate `json:"water"`
	Electricity PointRate `json:"electricity"`
}

func (o PointRateSaveWaterElectricity) PointRateSetting() PointRateSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointRateSetting(data)
}

// PointRatePublicTransport 公交地铁碳积分和积分兑换比例
type PointRatePublicTransport struct {
	Bus   PointRate `json:"bus"`
	Metro PointRate `json:"metro"`
}

func (o PointRatePublicTransport) PointRateSetting() PointRateSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointRateSetting(data)
}

// PointRateOnlineMeeting 线上会议碳积分和积分兑换比例
type PointRateOnlineMeeting struct {
	OneCity  PointRate `json:"oneCity"`  //异地会议
	ManyCity PointRate `json:"manyCity"` //同城会议
}

func (o PointRateOnlineMeeting) PointRateSetting() PointRateSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointRateSetting(data)
}

// PointRate 碳积分和积分兑换比例
type PointRate struct {
	CarbonCredit decimal.Decimal `json:"carbonCredit"`
	Point        int             `json:"point"`
}

func (o PointRate) PointRateSetting() PointRateSetting {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return PointRateSetting(data)
}

// Calc 根据碳积分计算可得积分
func (o PointRate) Calc(carbonCredit decimal.Decimal) int {
	return int(carbonCredit.Div(o.CarbonCredit).Mul(decimal.NewFromInt(int64(o.Point))).IntPart())
}
