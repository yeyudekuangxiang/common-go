package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

/*
收藏
*/

var DefaultCollectionController = CollectionController{}

type CollectionController struct {
}

func (ctr CollectionController) TopicCollection(c *gin.Context) (gin.H, error) {
	form := MyCollectionRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	collectionService := service.NewCollectionService(context.NewMioContext(context.WithContext(c.Request.Context())))

	collections, total, err := collectionService.TopicCollections(user.OpenId, form.Limit(), form.Offset())
	if err != nil {
		return nil, err
	}

	return gin.H{
		"list":     collections,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (ctr CollectionController) Collection(c *gin.Context) (gin.H, error) {
	form := CollectionRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)
	collectionService := service.NewCollectionService(ctx)
	err := collectionService.CollectionV2(form.ObjId, form.ObjType, user.OpenId)
	return nil, err
}

func (ctr CollectionController) CancelCollection(c *gin.Context) (gin.H, error) {
	form := CollectionRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)
	collectionService := service.NewCollectionService(context.NewMioContext(context.WithContext(c.Request.Context())))
	err := collectionService.CancelCollection(form.ObjId, form.ObjType, user.OpenId)
	return nil, err
}
