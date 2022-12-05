package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultFileExportController = FileExportController{}

type FileExportController struct {
}

func (ctr FileExportController) GetFileExportPageList(ctx *gin.Context) (gin.H, error) {
	var form GetFileExportPageListForm
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	t, _ := entity.FileExportTypeList.Find(int64(form.Type))
	list, total, err := service.DefaultFileExportService.GetPageList(repository.GetFileExportPageListBy{
		Type:           t,
		AdminId:        form.AdminId,
		Status:         form.Status,
		StartCreatedAt: model.Time{Time: form.StartCreatedAt},
		EndCreatedAt:   model.Time{Time: form.EndCreatedAt},
		Offset:         form.Offset(),
		Limit:          form.Limit(),
		OrderBy:        entity.OrderByList{entity.OrderByFileExportCreatedAtDesc},
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
func (ctr FileExportController) GetFileExportOptions(ctx *gin.Context) (interface{}, error) {
	data := service.DefaultFileExportService.GetFileExportStatusAndTypeList()
	return data, nil
}