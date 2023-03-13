package unittool

import (
	"github.com/shopspring/decimal"
)

type G float64

/*
数据记录、换算时，单位为克，精确到小数点后2位；
前端显示时，数据四舍五入后，精确到小数点后1位；
数据大于0但不足0.1的，按0.1显示；
10千克以内，前端使用单位“克”，如0.4 g，或9923 g；
10千克以上，前端使用单位“千克”，如10.4 kg，或 1982.3 kg；
10吨以上，前端使用单位“吨”，如 13.5吨
*/

func (k G) ToString() string {
	carbonDec := decimal.NewFromFloat(float64(k))
	if carbonDec.Cmp(decimal.NewFromFloat(10000)) == -1 {
		return carbonDec.Round(1).String() + "g"
	} else if carbonDec.Cmp(decimal.NewFromFloat(10000000)) == 1 {
		return carbonDec.Div(decimal.NewFromInt(1000000)).Round(1).String() + "T"
	} else {
		return carbonDec.Div(decimal.NewFromInt(1000)).Round(1).String() + "kg"
	}
}
