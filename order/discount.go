package order

type DiscountEnum int64

const (
	Point  DiscountEnum = 1 //积分
	Coupon DiscountEnum = 2 //优惠券
)

// 检测是否有使用过该优惠

func IsUseDiscount(discountType DiscountEnum, value int64) bool {
	return (value & (1 << (discountType - 1))) > 0
}

// 设置使用过的优惠

func SetDiscountValue(discountType DiscountEnum, value int64) int64 {
	return value | (1 << (discountType - 1))
}
