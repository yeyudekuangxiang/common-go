package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller"
	"mio/internal/app/mp2c/controller/api/api_types"
	entity2 "mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
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
	info, err := service.DefaultOrderService.SubmitOrderForEvent(srv_types.SubmitOrderForEventParam{
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

func (OrderController) SubmitOrderForEventGD(ctx *gin.Context) (gin.H, error) {
	form := api_types.SubmitOrderForEventGDForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	info, err := service.DefaultOrderService.SubmitOrderForEventGD(srv_types.SubmitOrderForEventGDParam{
		UserId:              user.ID,
		EventId:             form.EventId,
		OpenId:              user.OpenId,
		WechatServiceOpenId: form.WxServerOpenId,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"badgeInfo": info,
	}, nil
}

func (OrderController) GetUserOrderList(c *gin.Context) (interface{}, error) {
	page := controller.PageFrom{}
	if err := apiutil.BindForm(c, &page); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	list, total, err := service.DefaultOrderService.GetPageFullOrder(srv_types.GetPageFullOrderDTO{
		Openid:      user.OpenId,
		OrderSource: entity2.OrderSourceMio,
		Offset:      page.Offset(),
		Limit:       page.Limit(),
	})
	if err != nil {
		return nil, err
	}

	return controller.NewPageResult(list, total, page), nil
}
