package auth

import "testing"

import "fmt"

// 检测是否有使用过该优惠
func isUseDiscount(discountType int, value int) bool {
	return (value & (1 << (discountType - 1))) > 0
}

// 设置使用过的优惠
func setDiscountValue(discountType int, value int) int {
	return value | (1 << (discountType - 1))
}

func TestName2(t *testing.T) {
	value := 0                           // 初始存储值
	fmt.Println(isUseDiscount(1, value)) // 输出: false

	value = setDiscountValue(1, value)   // 设置使用过的优惠
	fmt.Println(isUseDiscount(1, value)) // 输出: true
}

func TestName(t *testing.T) {
	/*c := Client{
		UniTrustAppId: "71aec34abd4e44d091fc30e368d13bae",
		Token:         "92fa6eb418e51357a32615a17defde03",
	}
	auth, err := c.SendAuth(UserIdentityVerificationReq{
		Name:         "刘梅",
		IdentityCard: "211481199310056626",
		Phone:        "18840853003",
	})
	if err != nil {
		return
	}*/

	//	println(auth.Msg)
}
