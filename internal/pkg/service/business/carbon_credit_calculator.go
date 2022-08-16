package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model/entity/business"
	"time"
)

var DefaultCarbonCreditCalculatorService = CarbonCreditCalculatorService{}

// CarbonCreditCalculatorService 碳计算器
type CarbonCreditCalculatorService struct {
}

// CalcOnlineMeetingOneCity 根据同城会议时长计算获得多少碳积分
func (srv CarbonCreditCalculatorService) CalcOnlineMeetingOneCity(m time.Duration) decimal.Decimal {
	return decimal.NewFromFloat(m.Minutes() * 0.4).Round(2)
}

// CalcOnlineMeetingManyCity 根据多个城市会议时长计算获得多少碳积分
func (srv CarbonCreditCalculatorService) CalcOnlineMeetingManyCity(m time.Duration) decimal.Decimal {
	return decimal.NewFromFloat(m.Minutes() * 40).Round(2)
}

// CalcSaveWater 根绝节水量计算获得多少碳积分 水的单位升
func (srv CarbonCreditCalculatorService) CalcSaveWater(water int64) decimal.Decimal {
	panic("请配置我")
}

// CalcSaveElectricity 根绝节电量计算获得多少碳积分  电的单位度
func (srv CarbonCreditCalculatorService) CalcSaveElectricity(electricity int64) decimal.Decimal {
	panic("请配置我")
}

// CalcBus 乘坐公交车 km
func (srv CarbonCreditCalculatorService) CalcBus(bus float64) decimal.Decimal {
	return decimal.NewFromFloat(bus).Mul(decimal.NewFromFloat(111.45)).Round(2)
}

// CalcOEP 光盘行动
func (srv CarbonCreditCalculatorService) CalcOEP() decimal.Decimal {
	return decimal.NewFromFloat(111).Round(2)
}
func (srv CarbonCreditCalculatorService) CalcTrip(tripType business.TripType, distance decimal.Decimal) decimal.Decimal {
	return decimal.NewFromFloat(100).Round(2)
}

//CalcMetro 乘坐地铁 km
func (srv CarbonCreditCalculatorService) CalcMetro(metro float64) decimal.Decimal {
	return decimal.NewFromFloat(metro).Mul(decimal.NewFromFloat(134.05)).Round(2)
}

//CalcStep 步行	 km
func (srv CarbonCreditCalculatorService) CalcStep(step float64) decimal.Decimal {
	return decimal.NewFromFloat(step).Mul(decimal.NewFromFloat(134.05)).Round(2)
}

//CalcBike 骑行	 km
func (srv CarbonCreditCalculatorService) CalcBike(bike float64) decimal.Decimal {
	return decimal.NewFromFloat(bike).Mul(decimal.NewFromFloat(134.05)).Round(2)
}

/*
CalcEvCar 电车充电
electricity充电量 单位度
*/
func (srv CarbonCreditCalculatorService) CalcEvCar(electricity float64) decimal.Decimal {
	return decimal.NewFromFloat(electricity * 511).Round(2)
}
