package admin

import (
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultPointController = PointController{}

type PointController struct {
}

func (ctr PointController) GetPointRecordPageList(ctx *gin.Context) (gin.H, error) {
	var form GetPointRecordPageListFrom
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	pointTranService := service.NewPointTransactionService(context.NewMioContext(context.WithContext(ctx)))
	endTime := form.EndTime
	if !endTime.IsZero() {
		endTime = form.EndTime.Add(time.Hour*24 - time.Nanosecond)
	}

	endExpireTime := timetool.MustParse("2006-01-02", form.EndExpireTime)
	if !endExpireTime.IsZero() {
		endExpireTime = endExpireTime.EndOfDay()
	}

	list, total, err := pointTranService.PagePointRecord(service.GetPointTransactionPageListBy{
		UserId:          form.UserId,
		Nickname:        form.Nickname,
		OpenId:          form.OpenId,
		Phone:           form.Phone,
		StartTime:       model.Time{Time: form.StartTime},
		EndTime:         model.Time{Time: endTime},
		StartExpireTime: timetool.MustParse("2006-01-02", form.StartExpireTime).SqlNull(),
		EndExpireTime:   endExpireTime.SqlNull(),
		Type:            form.Type,
		Offset:          form.Offset(),
		Limit:           form.Limit(),
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
	pointTranService := service.NewPointTransactionService(context.NewMioContext(context.WithContext(ctx)))
	list := pointTranService.GetPointTransactionTypeList()
	return gin.H{
		"list": list,
	}, nil
}
func (ctr PointController) ExportPointRecordList(ctx *gin.Context) (gin.H, error) {
	var form ExportPointRecordListFrom
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	admin := apiutil.GetAuthAdmin(ctx)
	pointTranService := service.NewPointTransactionService(context.NewMioContext(context.WithContext(ctx)))
	err := pointTranService.ExportPointTransactionList(admin.ID, service.ExportPointTransactionListBy{
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
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	pointTranService := service.NewPointTransactionService(context.NewMioContext(context.WithContext(ctx)))
	list, total, err := pointTranService.PageAdjustPointRecord(service.GetPointAdjustRecordPageListParam{
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
	pointTranService := service.NewPointTransactionService(context.NewMioContext(context.WithContext(ctx)))
	list := pointTranService.GetAdjustPointTransactionTypeList()
	return gin.H{
		"list": list,
	}, nil
}
func (ctr PointController) AdjustUserPoint(ctx *gin.Context) (gin.H, error) {
	var form AdjustUserPointForm
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	admin := apiutil.GetAuthAdmin(ctx)

	pointService := service.NewPointService(context.NewMioContext(context.WithContext(ctx)))
	err := pointService.AdminAdjustUserPoint(admin.ID, service.AdminAdjustUserPointParam{
		OpenId: form.OpenId,
		Phone:  form.Phone,
		Type:   form.Type,
		Value:  form.Value,
		Note:   form.Note,
	})
	return nil, err
}
