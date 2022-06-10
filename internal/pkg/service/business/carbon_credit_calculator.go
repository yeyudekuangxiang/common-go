package business

import (
	"github.com/shopspring/decimal"
	"time"
)

var DefaultCarbonCreditCalculatorService = CarbonCreditCalculatorService{}

// CarbonCreditCalculatorService 碳计算器
type CarbonCreditCalculatorService struct {
}

// CalcOnlineMeetingOneCity 根据同城会议时长计算获得多少碳积分
func (srv CarbonCreditCalculatorService) CalcOnlineMeetingOneCity(m time.Duration) decimal.Decimal {
	return decimal.NewFromFloat(m.Minutes())
	panic("请配置我")
}

// CalcOnlineMeetingManyCity 根据多个城市会议时长计算获得多少碳积分
func (srv CarbonCreditCalculatorService) CalcOnlineMeetingManyCity(m time.Duration) decimal.Decimal {
	return decimal.NewFromFloat(m.Minutes() * 2)
	panic("请配置我")
}

// CalcSaveWater 根绝节水量计算获得多少碳积分 水的单位升
func (srv CarbonCreditCalculatorService) CalcSaveWater(water int64) decimal.Decimal {
	panic("请配置我")
}

// CalcSaveElectricity 根绝节电量计算获得多少碳积分  电的单位度
func (srv CarbonCreditCalculatorService) CalcSaveElectricity(electricity int64) decimal.Decimal {
	panic("请配置我")
}

// CalcBus 乘坐公交车
func (srv CarbonCreditCalculatorService) CalcBus(bus int64) decimal.Decimal {
	return decimal.NewFromInt(bus * 4)
	panic("请配置我")
}

//CalcMetro 乘坐地铁
func (srv CarbonCreditCalculatorService) CalcMetro(metro int64) decimal.Decimal {
	return decimal.NewFromInt(metro * 5)
	panic("请配置我")
}

/*
CalcEvCar 电车充电
electricity充电量 单位度
*/
func (srv CarbonCreditCalculatorService) CalcEvCar(electricity float64) decimal.Decimal {
	return decimal.NewFromFloat(electricity * 6)
	panic("请配置我")
}
