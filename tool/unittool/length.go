package unittool

import "github.com/shopspring/decimal"

type L int64

func (l L) ToString() string {
	var toString string
	if l >= 1000 {
		toString = decimal.NewFromInt(int64(l)).Div(decimal.NewFromInt(1000)).Round(2).String() + "km"
	} else {
		toString = decimal.NewFromInt(int64(l)).String() + "m"
	}
	return toString
}
