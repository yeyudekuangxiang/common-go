package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/kumiaoCommunity"
	"mio/internal/pkg/util/apiutil"
)

var DefaultTagController = TagController{}

type TagController struct {
}

func (ctr *TagController) List(c *gin.Context) (gin.H, error) {
	form := TagListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagAdminService := kumiaoCommunity.NewTagAdminService(ctx)

	list, total, err := tagAdminService.GetTagPageList(repository.GetTagPageListBy{
		Name:        form.Name,
		Description: form.Description,
		Offset:      form.Offset(),
		Limit:       form.Limit(),
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

func (ctr *TagController) Detail(c *gin.Context) (gin.H, error) {
	form := IDForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagAdminService := kumiaoCommunity.NewTagAdminService(ctx)

	topic := tagAdminService.Detail(form.ID)
	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr *TagController) Update(c *gin.Context) (gin.H, error) {
	form := UpdateTagRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	//更新帖子
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagAdminService := kumiaoCommunity.NewTagAdminService(ctx)

	err := tagAdminService.Update(repository.UpdateTag{
		ID: form.ID,
		CreateTag: repository.CreateTag{
			Name:        form.Name,
			Description: form.Description,
		},
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *TagController) Delete(c *gin.Context) (gin.H, error) {
	form := IDForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	//更新帖子
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagAdminService := kumiaoCommunity.NewTagAdminService(ctx)

	if err := tagAdminService.Delete(form.ID); err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *TagController) Create(c *gin.Context) (gin.H, error) {
	form := CreateTagRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagAdminService := kumiaoCommunity.NewTagAdminService(ctx)

	err := tagAdminService.Create(repository.CreateTag{
		Name:        form.Name,
		Description: form.Description,
		Image:       form.Image,
	})
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}
