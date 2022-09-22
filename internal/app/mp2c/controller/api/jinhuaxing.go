package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/platform"
	"mio/internal/pkg/util/apiutil"
)

var DefaultJinHuaXingController = JinHuaXingController{}

type JinHuaXingController struct {
}

func (ctr *JinHuaXingController) SendCoupon(ctx *gin.Context) (gin.H, error) {
	form := JinHuaXingForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//user := apiutil.GetAuthUser(ctx)
	service := platform.NewJinHuaXingService(context.NewMioContext())
	err := service.SendCoupon("123456789", form.Mobile)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
