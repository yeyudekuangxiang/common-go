package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (TopicController) List(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)

	list, total, err := service2.DefaultTopicService.GetTopicDetailPageList(repository.GetTopicPageListBy{
		ID:         form.ID,
		TopicTagId: form.TopicTagId,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
		UserId:     user.ID,
		Status:     entity.TopicStatusPublished,
	})
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0)
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	app.Logger.Infof("GetTopicDetailPageListByFlow user:%d form:%+v ids:%+v", user.ID, form, ids)

	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (TopicController) ListFlow(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)

	list, total, err := service2.DefaultTopicService.GetTopicDetailPageListByFlow(repository.GetTopicPageListBy{
		ID:         form.ID,
		TopicTagId: form.TopicTagId,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
		UserId:     user.ID,
	})
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0)
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	app.Logger.Infof("user:%d form:%+v ids:%+v", user.ID, form, ids)

	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}
func (TopicController) GetShareWeappQrCode(c *gin.Context) (gin.H, error) {
	form := GetWeappQrCodeFrom{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)

	buffers, contType, err := service2.DefaultTopicService.GetShareWeappQrCode(int(user.ID), form.TopicId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"buffers":     buffers,
		"contentType": contType,
	}, nil
}
func (TopicController) ChangeTopicLike(c *gin.Context) (gin.H, error) {
	form := ChangeTopicLikeForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)

	like, err := service2.TopicLikeService{}.ChangeLikeStatus(form.TopicId, int(user.ID))
	if err != nil {
		return nil, err
	}
	return gin.H{
		"status": like.Status,
	}, nil
}