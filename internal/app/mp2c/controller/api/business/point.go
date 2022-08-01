package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/business/businesstypes"
	"mio/internal/pkg/model/entity"
	ebusiness "mio/internal/pkg/model/entity/business"
	"mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultPointController = PointController{}

type PointController struct{}

func (PointController) GetPointRecordList(ctx *gin.Context) (gin.H, error) {
	form := businesstypes.GetPointRecordListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthBusinessUser(ctx)

	ptList := business.DefaultPointLogService.GetListBy(business.GetPointLogListParam{
		UserId:    user.ID,
		StartTime: form.Date,
		EndTime:   form.Date.AddDate(0, 1, 0).Add(-time.Nanosecond),
		OrderBy:   entity.OrderByList{ebusiness.OrderByPointLogCTDESC},
	})

	infoList := make([]businesstypes.PointLogInfo, 0)
	for _, pt := range ptList {
		fmt.Println(pt)
		infoList = append(infoList, businesstypes.PointLogInfo{
			ID:       pt.ID,
			Type:     pt.Type,
			TypeText: pt.Type.Text(),
			TimeStr:  pt.CreatedAt.Format("01.02 15:04:05"),
			Value:    pt.Value,
		})
	}
	return gin.H{
		"list": infoList,
	}, nil
}
