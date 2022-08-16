package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"time"
)

type CarbonController struct {
	service service.CarbonTransactionService
}

var DefaultCarbonController = CarbonController{
	service.NewCarbonTransactionService(context.NewMioContext())}

//测试

func (c CarbonController) Create(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionCreateForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//user := apiutil.GetAuthUser(ctx)
	//判断名称和图片是否存在.

	info, err := c.service.Create(api_types.CreateCarbonTransactionDto{
		OpenId:  form.OpenId,
		UserId:  form.UserId,
		Type:    form.Type,
		Value:   form.Value,
		Info:    fmt.Sprintf("{imageUrl=%s}", 1),
		AdminId: form.AdminId,
		Ip:      ctx.ClientIP(),
	})
	if err != nil {

	}
	return gin.H{
		"points": info,
	}, err
}

func (c CarbonController) Bank(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetCarbonTransactionBankForm{}
	user := apiutil.GetAuthUser(ctx)
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	list, total, err := c.service.Bank(api_types.GetCarbonTransactionBankDto{
		UserId: user.ID,
		Offset: form.Offset(),
		Limit:  form.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (c CarbonController) MyBank(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	bank, err := c.service.MyBank(api_types.GetCarbonTransactionMyBankDto{
		UserId: user.ID,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"myBank": bank,
	}, nil
}

func (c CarbonController) Info(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	info, err := c.service.Info(api_types.GetCarbonTransactionInfoDto{
		UserId: user.ID,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"info": info,
	}, nil
}

func (c CarbonController) Classify(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	ret, err := c.service.Classify(api_types.GetCarbonTransactionClassifyDto{
		UserId: user.ID,
	})
	if err != nil {
		return nil, err
	}
	//返回前整理
	var ListDateVo []string
	var ListValueVo []string

	total := decimal.NewFromFloat(ret.Total)
	for _, v := range ret.List {
		valDec := decimal.NewFromFloat(v.Val)
		overPer := "0%"
		if !total.IsZero() {
			overPer = valDec.Div(total).Round(2).Mul(decimal.NewFromInt(100)).String() + "%"
		}
		ListDateVo = append(ListDateVo, v.Key)
		ListValueVo = append(ListValueVo, overPer)
	}
	return gin.H{
		"dateList":  ListDateVo,
		"valueList": ListValueVo,
		"cover":     ret.Cover,
	}, nil
}

func (c CarbonController) History(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	//获取用户的
	list, err := c.service.History(api_types.GetCarbonTransactionHistoryDto{
		UserId:    user.ID,
		StartTime: time.Now().AddDate(0, 0, -14).Format("2006-01-02"),
		EndTime:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	})

	if err != nil {
		return nil, err
	}
	//整理 date为 key
	listMap := make(map[string]entity.CarbonTransactionDay)
	for _, j := range list {
		listMap[j.VDate.Format("2006-01-02")] = j
	}
	//整理最终列表 近2周的数据 i = 14
	var ListDateVo []string
	var ListValueVo []float64
	var ListValueStrVo []string

	for i := 1; i <= 14; i++ {
		day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		ListDateVo = append(ListDateVo, day)
		l, ok := listMap[day] //判断是否存在map集合中
		if ok {
			ListValueVo = append(ListValueVo, l.Value)
			ListValueStrVo = append(ListValueStrVo, util.CarbonToRate(l.Value))
		} else {
			ListValueVo = append(ListValueVo, 0)
			ListValueStrVo = append(ListValueStrVo, "0g")
		}
	}
	return gin.H{"dateList": ListDateVo, "valueList": ListValueVo, "valueListStr": ListValueStrVo}, nil
}

/****定时器*****/

func (c CarbonController) AddClassify(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	c.service.AddClassify(api_types.GetCarbonTransactionClassifyDto{
		UserId:    user.ID,
		StartTime: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
		EndTime:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	})
	return nil, nil
}

func (c CarbonController) AddHistory(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	c.service.AddHistory(api_types.GetCarbonTransactionClassifyDto{
		UserId:    user.ID,
		StartTime: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
		EndTime:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	})
	return nil, nil
}
