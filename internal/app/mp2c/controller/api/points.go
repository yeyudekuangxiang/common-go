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
	"mio/internal/pkg/service/platform"
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

	//查询openid是否是志愿汇id
	zyhService := platform.NewZyhService(context.NewMioContext())
	isVolunteer, _ := zyhService.CheckIsVolunteer(user.OpenId)
	recordInfoList := make([]api_types.PointRecordInfo, 0)
	for _, pt := range list {
		recordInfo := api_types.PointRecordInfo{}
		if err := util.MapTo(pt, &recordInfo); err != nil {
			return nil, err
		}
		recordInfo.TypeText = recordInfo.Type.Text() + getZyhTip(isVolunteer, recordInfo)
		recordInfo.TimeStr = recordInfo.CreateTime.Format("01-02 15:04:05")
		recordInfoList = append(recordInfoList, recordInfo)
	}

	return gin.H{
		"list": recordInfoList,
	}, nil
}

//志愿汇类型

var ZyhPointCollect = map[entity.PointTransactionType]string{
	entity.POINT_STEP:                    "步行",
	entity.POINT_COFFEE_CUP:              "自带咖啡杯",
	entity.POINT_BIKE_RIDE:               "骑行",
	entity.POINT_ECAR:                    "答题活动",
	entity.POINT_QUIZ:                    "答题活动",
	entity.POINT_JHX:                     "金华行",
	entity.POINT_POWER_REPLACE:           "电车换电",
	entity.POINT_DUIBA_INTEGRAL_RECHARGE: "兑吧虚拟商品充值积分",
	entity.POINT_RECYCLING_CLOTHING:      "旧物回收 oola衣物鞋帽",
	entity.POINT_RECYCLING_DIGITAL:       "旧物回收 oola数码",
	entity.POINT_RECYCLING_APPLIANCE:     "旧物回收 oola家电",
	entity.POINT_RECYCLING_BOOK:          "旧物回收 oola书籍",
	entity.POINT_FMY_RECYCLING_CLOTHING:  "旧物回收 fmy衣物鞋帽",
	entity.POINT_FAST_ELECTRICITY:        "快电",
	entity.POINT_REDUCE_PLASTIC:          "环保减塑",
	entity.POINT_CYCLING:                 "骑行",
}

func getZyhTip(isVolunteer bool, recordInfo api_types.PointRecordInfo) string {
	zyhTip := ""
	if isVolunteer {
		//2022-10-13 12:05:58 之前不提醒
		if recordInfo.CreateTime.Unix() > 1665633958 && recordInfo.Value > 0 {
			_, zyhOk := ZyhPointCollect[recordInfo.Type]
			if !zyhOk {
				zyhTip = "(该场景不转化志愿汇能源)"
			}
		}
		if recordInfo.Type == entity.POINT_QUIZ && recordInfo.Value == 2500 {
			zyhTip = "(该场景不转化志愿汇能源)"
		}
	}
	return zyhTip
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
