package order

import "github.com/shopspring/decimal"

// 检测是否有使用过该优惠

func PointToMoneyFen(point uint64, rate int64) int64 {
	if point == 0 || rate == 0 {
		return 0
	}
	pointDec := decimal.NewFromInt(int64(point))
	return pointDec.Div(decimal.NewFromInt(rate).Round(2)).Mul(decimal.NewFromInt(100)).BigInt().Int64()
}

// 设置使用过的优惠

func PointToMoneyYuan(point uint64, rate int64) float64 {
	moneyFen := PointToMoneyFen(point, rate)
	money, _ := decimal.NewFromInt(moneyFen).Div(decimal.NewFromInt(100)).Round(2).Float64()
	return money
}

func GetPointIntPart(point int64, rate int64) int64 {
	unit := decimal.NewFromInt(rate).Div(decimal.NewFromInt(100)).IntPart() //1分等于5积分，最少应用5积分
	intPart := decimal.NewFromInt(point).Div(decimal.NewFromInt(unit)).IntPart()
	return decimal.NewFromInt(intPart).Mul(decimal.NewFromInt(unit)).IntPart()
}
