package api

import (
	"github.com/gin-gonic/gin"
	entity2 "mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultTagController = TagController{}

type TagController struct {
}

func (TagController) List(c *gin.Context) (gin.H, error) {
	form := GetTagForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total, err := service.DefaultTagService.GetTagPageList(repository.GetTagPageListBy{
		ID:      form.ID,
		Offset:  form.Offset(),
		Limit:   form.Limit(),
		OrderBy: entity2.OrderByList{entity2.OrderByTagSortDesc},
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
