package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"time"
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
		EndTime:   model.Time{Time: form.EndTime.Add(time.Hour*24 - time.Nanosecond)},
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
	var form ExportPointRecordListFrom
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	admin := util.GetAuthAdmin(ctx)
	err := service.DefaultPointTransactionService.ExportPointTransactionList(admin.ID, service.ExportPointTransactionListBy{
		UserId:    form.UserId,
		Nickname:  form.Nickname,
		OpenId:    form.OpenId,
		Phone:     form.Phone,
		StartTime: model.Time{Time: form.StartTime},
		EndTime:   model.Time{Time: form.EndTime.Add(time.Hour*24 - time.Nanosecond)},
		Type:      form.Type,
	})
	return nil, err
}
func (ctr PointController) GetAdjustRecordPageList(ctx *gin.Context) (gin.H, error) {
	var form GetAdjustRecordPageListForm
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	list, total, err := service.DefaultPointTransactionService.GetAdjustRecordPageList(service.GetPointAdjustRecordPageListParam{
		OpenId:    form.OpenId,
		Phone:     form.Phone,
		Type:      form.Type,
		AdminId:   form.AdminId,
		Nickname:  form.Nickname,
		UserId:    form.UserId,
		StartTime: form.StartTime,
		EndTime:   form.EndTime.Add(time.Hour*24 - time.Nanosecond),
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
func (ctr PointController) GetAdjustPointTransactionTypeList(ctx *gin.Context) (gin.H, error) {
	list := service.DefaultPointTransactionService.GetAdjustPointTransactionTypeList()
	return gin.H{
		"list": list,
	}, nil
}
func (ctr PointController) AdjustUserPoint(ctx *gin.Context) (gin.H, error) {
	var form AdjustUserPointForm
	if err := util.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	admin := util.GetAuthAdmin(ctx)

	err := service.DefaultPointTransactionService.AdminAdjustUserPoint(admin.ID, service.AdminAdjustUserPointParam{
		OpenId: form.OpenId,
		Phone:  form.Phone,
		Type:   form.Type,
		Value:  form.Value,
		Note:   form.Note,
	})
	return nil, err
}
