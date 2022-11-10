package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultIndexIconController = NewIndexIconController()

func NewIndexIconController() IndexIconController {
	return IndexIconController{}
}

type IndexIconController struct {
}

//Create 创建icon
func (ctl IndexIconController) Create(c *gin.Context) (gin.H, error) {
	var form CreateIndexIconForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.NewIndexIconService(context.NewMioContext()).Create(entity.IndexIcon{
		Title:  form.Title,
		RowNum: form.RowNum,
		Sort:   form.Sort,
		Status: form.Status,
		IsOpen: form.IsOpen,
		Pic:    form.Pic,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

//Update 创建icon
func (ctl IndexIconController) Update(c *gin.Context) (gin.H, error) {
	var form UpdateIndexIconForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.NewIndexIconService(context.NewMioContext()).Update(entity.IndexIcon{
		ID:     form.Id,
		Title:  form.Title,
		RowNum: form.RowNum,
		Sort:   form.Sort,
		Status: form.Status,
		IsOpen: form.IsOpen,
		Pic:    form.Pic,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

//Delete 创建icon
func (ctl IndexIconController) Delete(c *gin.Context) (gin.H, error) {
	var form DeleteIndexIconForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.NewIndexIconService(context.NewMioContext()).DeleteById(repotypes.DeleteIndexIconDO{
		Id:     form.Id,
		Status: entity.IndexIconStatusDown,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (ctl IndexIconController) Page(c *gin.Context) (gin.H, error) {
	var form GetIndexIconPageForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total, err := service.NewIndexIconService(context.NewMioContext()).Page(repotypes.GetIndexIconPageDO{
		Offset: form.Offset(),
		Limit:  form.Limit(),
		Title:  form.Title,
		IsOpen: form.IsOpen,
	})

	voList := make([]api_types.IndexIconVO, 0)
	for _, activity := range list {
		voList = append(voList, api_types.IndexIconVO{
			ID:        activity.ID,
			Title:     activity.Title,
			Type:      activity.Type,
			RowNum:    activity.RowNum,
			Sort:      activity.Sort,
			Status:    activity.Status,
			IsOpen:    activity.IsOpen,
			Pic:       activity.Pic,
			CreatedAt: activity.CreatedAt.Format("2006.01.02 15:04:05"),
			UpdatedAt: activity.UpdatedAt.Format("2006.01.02 15:04:05"),
		})
	}
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
