package community

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	communityRepository "mio/internal/pkg/repository/community"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCommunityActivitiesTagController = ActivitiesTagController{}

type ActivitiesTagController struct {
}

// List 获取话题列表
func (ActivitiesTagController) List(c *gin.Context) (gin.H, error) {
	form := ActivitiesTagPageForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagService := community.NewCommunityActivitiesTagService(ctx)

	list, total, err := tagService.GetPageList(communityRepository.FindAllActivitiesTagParams{
		Offset: form.Offset(),
		Limit:  form.Limit(),
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

func (ActivitiesTagController) DetailTag(c *gin.Context) (gin.H, error) {
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagService := community.NewCommunityActivitiesTagService(ctx)

	tag, err := tagService.GetOne(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"tag": tag,
	}, nil
}
