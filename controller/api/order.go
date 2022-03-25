package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/service"
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
	order, err := service.DefaultOrderService.SubmitOrderForGreenMonday(service.SubmitOrderForGreenParam{
		AddressId: form.AddressId,
		UserId:    user.ID,
	})
	if err != nil {
		return nil, err
	}
	return order.ShortOrder(), nil
}
