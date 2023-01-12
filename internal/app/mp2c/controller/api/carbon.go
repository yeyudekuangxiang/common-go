package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	platformSrv "mio/internal/pkg/service/platform"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/encrypt"
	"sort"
	"strconv"
	"time"
)

type CarbonController struct {
	service service.CarbonTransactionService
}

var DefaultCarbonController = CarbonController{
	service.NewCarbonTransactionService(context.NewMioContext())}

func (c CarbonController) PointToCarbon(ctx *gin.Context) (gin.H, error) {
	c.service.PointToCarbon()
	return gin.H{}, nil
}

func GetSignToJava(params map[string]string) string {
	//排序
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + params[v]
	}
	return encrypt.Md5(signStr)
}

//Create 测试用
//{ "status":"ok", "errorMessage":" ", 'bizId':"20140730192133033", "carbonValue":22.22 }
func (c CarbonController) Create(ctx *gin.Context) (interface{}, error) {
	form := api_types.GetCarbonTransactionCreateForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return gin.H{
			"status":       "error",
			"errorMessage": err.Error(),
			"carbonValue":  "",
			"bizId":        "",
		}, nil
	}
	var carbonType entity.CarbonTransactionType
	var pointType entity.PointTransactionType
	switch form.CarbonType {
	case "Cycling":
		carbonType = entity.CARBON_CYCLING
		pointType = entity.POINT_CYCLING
	}
	if carbonType == "" {
		return gin.H{
			"status":       "error",
			"errorMessage": "类型有误",
			"carbonValue":  "",
			"bizId":        "",
		}, nil
	}
	params := make(map[string]string)
	params["serialNumber"] = form.SerialNumber
	params["carbonType"] = form.CarbonType
	params["carbonValue"] = form.CarbonValue
	params["pointValue"] = form.PointValue
	params["uid"] = form.Uid
	params["time"] = form.Time
	params["info"] = form.Info
	params["privateKey"] = "mio2022"
	params["info"] = form.Info

	sign := GetSignToJava(params)
	if sign != form.Sign {
		return gin.H{
			"status":       "error",
			"errorMessage": "验证有误",
			"carbonValue":  "",
			"bizId":        "",
		}, nil
	}

	carbonValue, carbonErr := strconv.ParseFloat(form.CarbonValue, 64)
	if carbonErr != nil {
		return gin.H{
			"status":       "error",
			"errorMessage": carbonErr.Error(),
			"carbonValue":  "",
			"bizId":        "",
		}, nil
	}

	info, err := c.service.Create(api_types.CreateCarbonTransactionDto{
		OpenId:   form.Uid,
		Type:     carbonType,
		AddValue: carbonValue,
		Info:     form.Info,
		AdminId:  0,
		Ip:       "", //ctx.ClientIP()*/
	})
	if err != nil {
		return gin.H{
			"status":       "error",
			"errorMessage": err.Error(),
			"carbonValue":  "",
			"bizId":        "",
		}, nil
	}

	point, errPoint := strconv.ParseInt(form.PointValue, 10, 64)
	if errPoint != nil {
		return gin.H{
			"status":       "error",
			"errorMessage": errPoint.Error(),
			"carbonValue":  "",
			"bizId":        form.SerialNumber,
		}, nil
	}

	//同步到志愿汇
	if point >= 0 {
		sendType := "0"
		serviceZyh := platformSrv.NewZyhService(context.NewMioContext())
		messageCode, messageErr := serviceZyh.SendPoint(sendType, form.Uid, strconv.FormatInt(point, 10))
		if messageCode != "30000" {
			//发送结果记录到日志
			msgErr := ""
			if messageErr != nil {
				msgErr = messageErr.Error()
			}
			serviceZyh.CreateLog(srv_types.GetZyhLogAddDTO{
				Openid:         form.Uid,
				PointType:      pointType,
				Value:          point,
				ResultCode:     messageCode,
				AdditionalInfo: msgErr,
				TransactionId:  form.SerialNumber,
			})
		}
	}

	return gin.H{
		"status":       "ok",
		"errorMessage": "",
		"bizId":        "",
		"carbonValue":  info,
	}, nil
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
	c.service.AddClassifyByUid(user.ID)
	ret, err := c.service.Classify(api_types.GetCarbonTransactionClassifyDto{
		UserId: user.ID,
	})
	if err != nil {
		return gin.H{}, err
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
		listMap[j.VDate.Format("01-02")] = j
	}

	//整理最终列表 近2周的数据 i = 14
	var ListDateVo []string
	var ListValueVo []float64
	var ListValueStrVo []string

	for i := 13; i >= 1; i-- {
		day := time.Now().AddDate(0, 0, -i).Format("01-02")
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
	//今天的数据加入进去
	carbonToday := c.service.GetTodayCarbon(user.ID)
	ListDateVo = append(ListDateVo, time.Now().Format("01-02"))
	ListValueVo = append(ListValueVo, carbonToday)
	ListValueStrVo = append(ListValueStrVo, util.CarbonToRate(carbonToday))
	return gin.H{"dateList": ListDateVo, "valueList": ListValueVo, "valueListStr": ListValueStrVo}, nil
}

/****定时器*****/

func (c CarbonController) AddClassify(ctx *gin.Context) (gin.H, error) {
	c.service.AddClassify()
	return nil, nil
}

func (c CarbonController) AddHistory(ctx *gin.Context) (gin.H, error) {
	c.service.AddHistory()
	return nil, nil
}
