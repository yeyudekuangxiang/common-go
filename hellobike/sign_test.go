package hellobike

import (
	"fmt"
	"testing"
)

//{"orderNo":"TP20230317181438200002103590028"}

func TestName(t *testing.T) {
	u := (*User)(nil)
	u.Name()
	u.Age()

}

type User struct {
}

func (u *User) Name() {
	fmt.Println("1111")
}
func (u *User) Age() {
	fmt.Println("2222")
}

func TestUser(t *testing.T) {
	u := User{}
	u.Name()
	u.Age()
}

func TestRefundCardParam(t *testing.T) {

}

func TestRefundCard(t *testing.T) {
	c := Client{
		AppId:  "20230302145050102",
		AppKey: "d9244321dc3246caa54a29e7c156dd0c",
		Domain: "https://openapi.hellobike.com/open/api",
	}
	resp, _ := c.RefundCard(RefundCardParam{
		ActivityId:    "H3979885952972083867",
		OrderNo:       "TP20230317181438200002103590028",
		MobilePhone:   "18840853003",
		TransactionId: "",
	})
	println(resp.Code)

}

func TestBikeCard(t *testing.T) {
	c := Client{
		AppId:  "20230302145050102",
		AppKey: "d9244321dc3246caa54a29e7c156dd0c",
		Domain: "https://openapi.hellobike.com/bike/activity",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, _ := c.SendCoupon(SendCouponParam{
		ActivityId:    "H3979885952972083867",
		MobilePhone:   "18840853003",
		TransactionId: "202009091447R0001243",
	})
	println(resp.ErrorCode)

}

/*func TestHelloBike(t *testing.T) {
	c := Client{
		AppId:  "20200907153742407",
		AppKey: "75e3747b359246379b2447dfd5090b8a",
	}
	resp, _ := c.SendCoupon(SendCouponParam{
		ActivityId:    "H3979885952972083867",
		MobilePhone:   "13661502232",
		TransactionId: "202009091447R0001243",
	})
	println(resp.ErrorCode)
}
*/
