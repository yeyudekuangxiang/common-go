package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/service_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultOrderController = OrderController{}

type OrderController struct {
}

func (OrderController) SubmitOrderForGreen(ctx *gin.Context) (interface{}, error) {
	form := SubmitOrderForGreenForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	order, err := service.DefaultOrderService.SubmitOrderForGreenMonday(service.SubmitOrderForGreenParam{
		AddressId: form.AddressId,
		UserId:    user.ID,
	})
	if err != nil {
		return nil, err
	}
	return order.ShortOrder(), nil
}
func (OrderController) SubmitOrderForEvent(ctx *gin.Context) (gin.H, error) {
	form := api_types.SubmitOrderForEventForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	info, err := service.DefaultOrderService.SubmitOrderForEvent(service_types.SubmitOrderForEventParam{
		UserId:  user.ID,
		EventId: form.EventId,
	})

	if err != nil {
		return nil, err
	}
	return gin.H{
		"badgeInfo": info,
	}, nil
}
