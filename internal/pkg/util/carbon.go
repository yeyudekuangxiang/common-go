package util

import "github.com/shopspring/decimal"

/*
树木碳吸引量换算：
按1棵成年的沙棘树每年吸收二氧化碳 1.66 kg 来计算；
若用户碳减排量为A千克，则相当于B（B = A / 1.66 kg）棵沙棘树1年碳吸收量；
若B<1棵，则计算 相当于1棵沙棘树C（C = A / (1.66 kg/365天)）天的碳吸收量；
*/

func CarbonToTree(carbon float64) (string, string) {
	rateDec := decimal.NewFromFloat(1660) //转化率
	carbonDec := decimal.NewFromFloat(carbon)
	if carbonDec.Cmp(rateDec) == -1 {
		//大于等于一棵树
		oneDayCarbon := rateDec.Div(decimal.NewFromInt(365))

		val := carbonDec.Div(oneDayCarbon).Round(1).String()
		msg := "棵沙棘树," + val + "天的碳吸收量"
		return "1", msg
	} else {
		//不够一棵树
		val := carbonDec.Div(rateDec).Round(1).String()
		msg := "棵沙棘树,1年碳吸收量"
		return val, msg
	}
}

/*
数据记录、换算时，单位为克，精确到小数点后2位；
前端显示时，数据四舍五入后，精确到小数点后1位；
数据大于0但不足0.1的，按0.1显示；
10千克以内，前端使用单位“克”，如0.4 g，或9923 g；
10千克以上，前端使用单位“千克”，如10.4 kg，或 1982.3 kg；
10吨以上，前端使用单位“吨”，如 13.5吨
*/

func CarbonToRate(carbon float64) string {
	carbonDec := decimal.NewFromFloat(carbon)
	if carbonDec.Cmp(decimal.NewFromFloat(10000)) == -1 {
		return carbonDec.Round(1).String() + "g"
	} else if carbonDec.Cmp(decimal.NewFromFloat(10000000)) == 1 {
		return carbonDec.Div(decimal.NewFromInt(1000000)).Round(1).String() + "t"
	} else {
		return carbonDec.Div(decimal.NewFromInt(1000)).Round(1).String() + "kg"
	}
}
