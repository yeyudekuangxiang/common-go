package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/converttool"
	"mio/internal/pkg/core/context"
	entity2 "mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/kumiaoCommunity"
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
	tagService := kumiaoCommunity.NewTagService(ctx)

	list, total, err := tagService.GetTagPageList(repository.GetTagPageListBy{
		ID:      form.ID,
		Offset:  form.Offset(),
		Limit:   form.Limit(),
		OrderBy: entity2.OrderByList{entity2.OrderByTagSortDesc},
		Status:  converttool.PointerInt(1),
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
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagService := kumiaoCommunity.NewTagService(ctx)

	tag, err := tagService.GetOne(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": tag,
	}, nil
}
