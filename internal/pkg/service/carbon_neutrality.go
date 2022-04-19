package service

import "math"

const conversionRate float64 = 1.0

var DefaultCarbonNeutralityService = CarbonNeutralityService{}

type CarbonNeutralityService struct {
}

func (srv CarbonNeutralityService) calculateCO2ByStep(steps int) int {
	return int(math.Floor(math.Floor(float64(steps)/1000) * conversionRate))
}
