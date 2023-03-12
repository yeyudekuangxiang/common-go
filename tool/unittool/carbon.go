package unittool

import "github.com/shopspring/decimal"

func CarbonToRate(carbon float64) string {
	carbonDec := decimal.NewFromFloat(carbon)
	if carbonDec.Cmp(decimal.NewFromFloat(10000)) == -1 {
		return carbonDec.Round(1).String() + "g"
	} else if carbonDec.Cmp(decimal.NewFromFloat(10000000)) == 1 {
		return carbonDec.Div(decimal.NewFromInt(1000000)).Round(1).String() + "T"
	} else {
		return carbonDec.Div(decimal.NewFromInt(1000)).Round(1).String() + "kg"
	}
}
