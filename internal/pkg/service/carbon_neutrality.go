package service

import "math"

const conversionRate float64 = 21.45 / 1000

var DefaultCarbonNeutralityService = CarbonNeutralityService{}

type CarbonNeutralityService struct {
}

//保留结果保留两位小数
func (srv CarbonNeutralityService) calculateCO2ByStep(steps int64) float64 {
	co2 := float64(steps) * conversionRate
	return math.Round(co2*100) / 100
}
