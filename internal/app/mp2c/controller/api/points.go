package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultPointController = PointsController{}

type PointsController struct {
}

func (PointsController) GetPointTransactionList(ctx *gin.Context) (gin.H, error) {
	form := GetPointTransactionListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	endTime := form.EndTime
	if !endTime.IsZero() {
		endTime = endTime.Add(time.Hour * 24).Add(-time.Nanosecond)
	}
	user := apiutil.GetAuthUser(ctx)
	pointTranService := service.NewPointTransactionService(context.NewMioContext(context.WithContext(ctx)))
	list := pointTranService.GetListBy(repository.GetPointTransactionListBy{
		OpenId:    user.OpenId,
		StartTime: model.Time{Time: form.StartTime},
		EndTime:   model.Time{Time: endTime},
		OrderBy:   entity.OrderByList{entity.OrderByPointTranCTDESC},
	})

	recordInfoList := make([]api_types.PointRecordInfo, 0)
	for _, pt := range list {
		recordInfo := api_types.PointRecordInfo{}
		if err := util.MapTo(pt, &recordInfo); err != nil {
			return nil, err
		}
		recordInfo.TypeText = recordInfo.Type.Text()
		recordInfo.TimeStr = recordInfo.CreateTime.Format("01-02 15:04:05")
		recordInfoList = append(recordInfoList, recordInfo)
	}

	return gin.H{
		"list": recordInfoList,
	}, nil
}
func (PointsController) GetPoint(ctx *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(ctx)
	pointService := service.NewPointService(context.NewMioContext())
	point, err := pointService.FindByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"points": point.Balance,
	}, nil
}

func (PointsController) MyReward(ctx *gin.Context) (gin.H, error) {
	form := controller.PageFrom{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	//我的奖励分类
	var myRewardType []entity.PointTransactionType
	myRewardType = append(myRewardType, entity.POINT_ARTICLE, entity.POINT_LIKE, entity.POINT_RECOMMEND, entity.POINT_COMMENT)

	pointTranactionService := service.NewPointTransactionService(context.NewMioContext())
	record, total, err := pointTranactionService.PagePointRecord(service.GetPointTransactionPageListBy{
		OpenId: user.OpenId,
		Types:  myRewardType,
		Offset: form.Offset(),
		Limit:  form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     record,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}
