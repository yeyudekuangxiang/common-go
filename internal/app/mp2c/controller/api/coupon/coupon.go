package coupon

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
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
func (CouponController) GetPageUserCouponRecord(c *gin.Context) (interface{}, error) {
	form := controller.PageFrom{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	list, total, err := service.DefaultCouponService.GetPageUserCouponRecord(srv_types.GetPageCouponRecordDTO{
		OpenId: user.OpenId,
		Offset: form.Offset(),
		Limit:  form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	return controller.NewPageResult(list, total, form), nil
}
func (CouponController) RedeemCode(c *gin.Context) (interface{}, error) {
	form := RedeemCodeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	result, err := service.DefaultCouponService.RedeemCouponFromCode(service.RedeemCouponByCodeParam{
		OpenId:       user.OpenId,
		RedeemCodeId: form.RedeemCodeId,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
