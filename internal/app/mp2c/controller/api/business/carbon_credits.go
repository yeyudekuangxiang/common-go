package business

import (
	"github.com/gin-gonic/gin"
	business2 "mio/internal/pkg/repository/business"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultCarbonCreditsController = CarbonCreditsController{}

type CarbonCreditsController struct{}

func (CarbonCreditsController) GetCarbonCreditLogInfoList(ctx *gin.Context) (gin.H, error) {
	form := GetCarbonCreditLogInfoListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthBusinessUser(ctx)

	infoList := business.DefaultCarbonCreditsLogService.GetCarbonCreditLogInfoList(business.GetCarbonCreditLogInfoListParam{
		UserId:    user.ID,
		StartTime: form.Date,
		EndTime:   form.Date.AddDate(0, 1, 0).Add(-time.Nanosecond),
	})
	return gin.H{
		"list": infoList,
	}, nil
}

func (CarbonCreditsController) GetCarbonCreditLogSortedList(ctx *gin.Context) (gin.H, error) {
	form := GetCarbonCreditLogSortedListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	//查询减碳排行
	list := business.DefaultCarbonCreditsLogService.GetCarbonCreditLogSortedList(business.GetCarbonCreditLogSortedListParam{
		StartTime: form.StartTime,
		EndTime:   form.EndTime.Add(time.Hour*24 - time.Nanosecond),
	})
	res := business.DefaultCarbonCreditsLogService.FormatCarbonCreditLogList(list)
	return gin.H{
		"list": res,
	}, nil
}

func (CarbonCreditsController) GetUserCarbonCreditLogSortedList(ctx *gin.Context) (gin.H, error) {
	form := GetCarbonCreditLogSortedListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthBusinessUser(ctx)
	//查询减碳排行
	list := business.DefaultCarbonCreditsLogService.GetCarbonCreditLogSortedList(business.GetCarbonCreditLogSortedListParam{
		StartTime: form.StartTime,
		EndTime:   form.EndTime.Add(time.Hour*24 - time.Nanosecond),
		UserId:    user.ID,
	})
	res := business.DefaultCarbonCreditsLogService.FormatCarbonCreditLogList(list)
	return gin.H{
		"list": res,
	}, nil
}

func (CarbonCreditsController) GetCarbonCreditLogSortedListHistory(ctx *gin.Context) (gin.H, error) {
	list := business.DefaultCarbonCreditsLogService.GetCarbonCreditLogListHistoryBy(business2.GetCarbonCreditsLogSortedListBy{})
	return gin.H{
		"list": list,
	}, nil
}

func (CarbonCreditsController) GetCarbonCreditLogSortedListUserHistory(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthBusinessUser(ctx)
	list := business.DefaultCarbonCreditsLogService.GetCarbonCreditLogListHistoryBy(business2.GetCarbonCreditsLogSortedListBy{UserId: user.ID})
	return gin.H{
		"list": list,
	}, nil
}
