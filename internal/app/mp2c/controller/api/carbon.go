package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCarbonController = CarbonController{}

type CarbonController struct {
}

func (c CarbonController) Create(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionBankForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//user := apiutil.GetAuthUser(ctx)
	//判断名称和图片是否存在
	a, err := service.DefaultCarbonTransactionService.Create(api_types.CreateCarbonTransactionDto{
		OpenId:  "1",
		UserId:  5,
		Type:    entity.CARBON_COFFEE_CUP,
		Value:   1,
		Info:    fmt.Sprintf("{imageUrl=%s}", 1),
		AdminId: 1,
	})
	if err != nil {

	}
	return gin.H{
		"points": a,
	}, err
}
func (c CarbonController) Bank(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionBankForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	a, _ := service.DefaultCarbonTransactionService.Bank(api_types.GetCarbonTransactionBankDto{
		UserId: 5,
	})
	return gin.H{
		"points": a,
	}, nil
}

func (c CarbonController) MyBank(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionMyBankForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	a, _ := service.DefaultCarbonTransactionService.MyBank(api_types.GetCarbonTransactionMyBankDto{
		UserId: 5,
	})
	return gin.H{
		"points": a,
	}, nil
}

func (c CarbonController) Info(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionInfoForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	service.DefaultCarbonTransactionService.Info(api_types.GetCarbonTransactionInfoDto{
		UserId: 5,
	})
	return gin.H{
		"points": 1,
	}, nil
}

func (c CarbonController) Classify(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionClassifyForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	service.DefaultCarbonTransactionService.Classify(api_types.GetCarbonTransactionClassifyDto{
		UserId: 5,
	})
	return gin.H{
		"points": form,
	}, nil
}

func (c CarbonController) History(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionHistoryForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	service.DefaultCarbonTransactionService.History(api_types.GetCarbonTransactionHistoryDto{
		UserId: 5,
	})
	return gin.H{
		"points": form,
	}, nil
}
