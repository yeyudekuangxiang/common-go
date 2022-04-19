package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

var DefaultPointController = PointsController{}

type PointsController struct {
}

func (PointsController) GetPointTransactionList(ctx *gin.Context) (gin.H, error) {
	form := GetPointTransactionListForm{}
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(ctx)
	list := service.DefaultPointTransactionService.GetListBy(repository.GetPointTransactionListBy{
		OpenId:    user.OpenId,
		StartTime: model.Time{Time: form.StartTime},
		EndTime:   model.Time{Time: form.EndTime},
		OrderBy:   entity.OrderByList{entity.OrderByPointTranCTDESC},
	})

	return gin.H{
		"list": list,
	}, nil
}
func (PointsController) GetPoint(ctx *gin.Context) (gin.H, error) {
	user := util.GetAuthUser(ctx)
	point, err := service.DefaultPointService.FindByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"points": point.Balance,
	}, nil
}
