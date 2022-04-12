package api

import (
	"github.com/gin-gonic/gin"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

var DefaultOrderController = OrderController{}

type OrderController struct {
}

func (OrderController) SubmitOrderForGreen(ctx *gin.Context) (interface{}, error) {
	form := SubmitOrderForGreenForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := util.GetAuthUser(ctx)
	order, err := service2.DefaultOrderService.SubmitOrderForGreenMonday(service2.SubmitOrderForGreenParam{
		AddressId: form.AddressId,
		UserId:    user.ID,
	})
	if err != nil {
		return nil, err
	}
	return order.ShortOrder(), nil
}