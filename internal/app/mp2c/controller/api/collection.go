package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
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
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	collectionService := service.NewCollectionService(ctx)

	collections, total, err := collectionService.TopicCollections(user.OpenId, form.Limit(), form.Offset())
	if err != nil {
		return nil, err
	}

	resList := make([]*entity.TopicItemRes, 0)

	//点赞数据
	likeMap := make(map[int64]struct{}, 0)
	topicLikeService := service.NewTopicLikeService(ctx)
	likeList, _ := topicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = struct{}{}
		}
	}

	//评论数据
	ids := make([]int64, 0) //topicId
	for _, item := range collections {
		ids = append(ids, item.Id)
	}

	rootCommentCount := service.DefaultTopicService.GetRootCommentCount(ids)

	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}

	//组装数据---点赞数据 收藏数据
	for _, item := range collections {
		res := item.TopicItemRes()
		//评论数量
		res.CommentCount = topic2comment[res.Id]
		//是否点赞
		if _, ok := likeMap[res.Id]; ok {
			res.IsLike = 1
		}
		//是否收藏
		res.IsCollection = 1
		resList = append(resList, res)
	}

	return gin.H{
		"list":     resList,
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
