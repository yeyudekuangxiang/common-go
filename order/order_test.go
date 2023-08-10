package order

import (
	"fmt"
	"testing"
)

func TestDiscount(t *testing.T) {
	discountType := int64(0)
	discountType = SetDiscountValue(Coupon, discountType)
	isPoint := IsUseDiscount(Point, discountType)
	println(isPoint)
	isCoupon := IsUseDiscount(Coupon, discountType)
	println(isCoupon)
}

func TestUnit(t *testing.T) {
	a := PointToMoneyFen(100, 500)
	fmt.Println(a)

	b := PointToMoneyYuan(100, 500)
	fmt.Println(b)
}
