package order

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	a := GetPointIntPart(609, 500)

	println(a)
}
func TestDiscount(t *testing.T) {

	a := int64(-500)
	b := uint64(-a)
	println(b)
	discountType := int64(0)
	discountType = SetDiscountValue(Coupon, discountType)
	isPoint := IsUseDiscount(Point, discountType)
	println(isPoint)
	isCoupon := IsUseDiscount(Coupon, discountType)
	println(isCoupon)

}

func TestUnit(t *testing.T) {
	a := PointToMoneyFen(605, 500)
	fmt.Println(a)

	b := PointToMoneyYuan(605, 500)
	fmt.Println(b)
}
