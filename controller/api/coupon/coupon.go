package coupon

import (
	"github.com/gin-gonic/gin"
	"mio/service"
)

var DefaultCouponController = CouponController{}

type CouponController struct {
}

func (CouponController) CouponListOfOpenid(c *gin.Context) (gin.H, error) {
	openid := c.Query("openid")
	list, err := service.DefaultCouponService.CouponListOfOpenid(openid, []string{"80defb4f-f002-442f-b3a8-6c28a04ba814", "evcard0point"})

	return gin.H{
		"records": list,
	}, err
}
