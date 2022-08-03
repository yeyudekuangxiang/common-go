package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
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
	list, total, err := service.DefaultTagAdminService.GetTagPageList(repository.GetTagPageListBy{
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
	topic, err := service.DefaultTagAdminService.Detail(form.ID)
	if err != nil {
		return nil, err
	}
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
	err := service.DefaultTagAdminService.Update(repository.UpdateTag{
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
	if err := service.DefaultTagAdminService.Delete(form.ID); err != nil {
		return nil, err
	}
	return nil, nil
}
