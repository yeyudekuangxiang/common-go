package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/repository"
	"mio/service"
)

var DefaultTagController = TagController{}

type TagController struct {
}

func (TagController) List(c *gin.Context) (gin.H, error) {
	form := GetTagForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total, err := service.DefaultTagService.GetTagPageList(repository.GetTagPageListBy{
		ID:     form.ID,
		Offset: form.Offset(),
		Limit:  form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"data":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil

}
