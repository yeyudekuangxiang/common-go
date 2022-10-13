package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/star_charge"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
)

var DefaultStarChargeController = NewStarChargeController()

func NewStarChargeController() StarChargeController {
	return StarChargeController{}
}

type StarChargeController struct {
}

func (s StarChargeController) SendCoupon(ctx *gin.Context) (gin.H, error) {
	var form StarChargeSendCoupon
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//查用户
	user, err := service.DefaultUserService.GetUserByOpenId(form.OpenId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("用户不存在")
	}
	if user.PhoneNumber == "" {
		return nil, errno.ErrCommon.WithMessage("用户未绑定手机号")
	}
	starChargeService := star_charge.NewStarChargeService(context.NewMioContext())
	token, err := starChargeService.GetAccessToken()
	if err != nil {
		return nil, err
	}
	err = starChargeService.SendCoupon(user.OpenId, user.PhoneNumber, starChargeService.ProvideId, token)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
