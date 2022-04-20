package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

var DefaultPointController = PointController{}

type PointController struct {
}

func (ctr PointController) GetPointRecordPageList(ctx *gin.Context) (gin.H, error) {
	var form GetPointRecordPageListFrom
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	list, total, err := service.DefaultPointTransactionService.GetPageListBy(service.GetPointTransactionPageListBy{
		UserId:    form.UserId,
		Nickname:  form.Nickname,
		OpenId:    form.OpenId,
		Phone:     form.Phone,
		StartTime: model.Time{Time: form.StartTime},
		EndTime:   model.Time{Time: form.EndTime},
		Type:      form.Type,
		Offset:    form.Offset(),
		Limit:     form.Limit(),
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
func (ctr PointController) GetPointTypeList(ctx *gin.Context) (gin.H, error) {
	list := service.DefaultPointTransactionService.GetPointTransactionTypeList()
	return gin.H{
		"list": list,
	}, nil
}
func (ctr PointController) ExportPointRecordList(ctx *gin.Context) (gin.H, error) {
	var form GetPointRecordPageListFrom
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	err := service.DefaultPointTransactionService.ExportPointTransactionList(0, service.GetPointTransactionPageListBy{
		UserId:    form.UserId,
		Nickname:  form.Nickname,
		OpenId:    form.OpenId,
		Phone:     form.Phone,
		StartTime: model.Time{Time: form.StartTime},
		EndTime:   model.Time{Time: form.EndTime},
		Type:      form.Type,
		Offset:    form.Offset(),
		Limit:     form.Limit(),
	})
	return nil, err
}
