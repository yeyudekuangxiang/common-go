package unittool

import (
	"github.com/shopspring/decimal"
)

type GWeight float64

func (k GWeight) ToString() string {
	carbonDec := decimal.NewFromFloat(float64(k))
	if carbonDec.Cmp(decimal.NewFromFloat(10000)) == -1 {
		return carbonDec.Round(1).String() + "g"
	} else if carbonDec.Cmp(decimal.NewFromFloat(10000000)) == 1 {
		return carbonDec.Div(decimal.NewFromInt(1000000)).Round(1).String() + "T"
	} else {
		return carbonDec.Div(decimal.NewFromInt(1000)).Round(1).String() + "kg"
	}
}
