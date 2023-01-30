package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/jhx"
	"mio/internal/pkg/service/platform/star_charge"
	"mio/internal/pkg/service/platform/ytx"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
)

var DefaultPrivateController = NewPrivateController()

func NewPrivateController() PrivateController {
	return PrivateController{}
}

type PrivateController struct {
}

func (c PrivateController) SendCoupon(ctx *gin.Context) (gin.H, error) {
	var form SendCouponReq
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
	mio := context.NewMioContext()
	switch form.Ch {
	case "lvmiao":
		err = c.sendCouponForStarCharge(mio, *user)
		if err != nil {
			return nil, err
		}
	case "jhx":
		err = c.sendCouponForJhx(mio, *user)
		if err != nil {
			return nil, err
		}
	case "ytx":
		err = c.sendCouponForYtx(mio, form.Amount, *user)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func (c PrivateController) sendCouponForStarCharge(ctx *context.MioContext, user entity.User) error {
	starChargeService := star_charge.NewStarChargeService(ctx)
	token, err := starChargeService.GetAccessToken()
	if err != nil {
		return err
	}
	err = starChargeService.SendCoupon(user.OpenId, user.PhoneNumber, starChargeService.ProvideId, token)
	return err
}

func (c PrivateController) sendCouponForJhx(ctx *context.MioContext, user entity.User) error {
	jhxService := jhx.NewJhxService(ctx)
	_, err := jhxService.SendCoupon(1000, user)
	if err != nil {
		return err
	}
	return nil
}

func (c PrivateController) sendCouponForYtx(ctx *context.MioContext, amount float64, user entity.User) error {
	bdscene := service.DefaultBdSceneService.FindByCh("yitongxing")
	var options []ytx.Options
	options = append(options, ytx.WithPoolCode(bdscene.AppId2))
	options = append(options, ytx.WithSecret(bdscene.Secret))
	options = append(options, ytx.WithAppId(bdscene.AppId))
	options = append(options, ytx.WithDomain(bdscene.Domain))
	ytxService := ytx.NewYtxService(ctx, options...)
	if amount == 0 {
		amount = 1.00
	}
	_, err := ytxService.SendCoupon(1001, amount, user)
	if err != nil {
		return err
	}
	return nil
}
