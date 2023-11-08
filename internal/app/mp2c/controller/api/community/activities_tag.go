package community

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCommunityActivitiesTagController = ActivitiesTagController{}

type ActivitiesTagController struct {
}

// List 获取话题列表
func (ActivitiesTagController) List(c *gin.Context) (gin.H, error) {
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	tagService := community.NewCommunityActivitiesTagService(ctx)

	list, err := tagService.List()
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list": list,
	}, nil
}

func (ActivitiesTagController) DetailTag(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
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
