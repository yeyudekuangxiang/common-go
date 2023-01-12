package community

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	entity2 "mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/util/apiutil"
)

var DefaultTagController = TagController{}

type TagController struct {
}

// List 获取话题列表
func (TagController) List(c *gin.Context) (gin.H, error) {
	form := GetTagForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagService := community.NewTagService(ctx)
	
	status := 1 // will be of type int
	statusPtr := &status

	list, total, err := tagService.GetTagPageList(repository.GetTagPageListBy{
		ID:      form.ID,
		Offset:  form.Offset(),
		Limit:   form.Limit(),
		OrderBy: entity2.OrderByList{entity2.OrderByTagSortDesc},
		Status:  statusPtr,
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

func (TagController) DetailTag(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagService := community.NewTagService(ctx)

	tag, err := tagService.GetOne(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": tag,
	}, nil
}
