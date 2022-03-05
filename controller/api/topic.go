package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/util"
	"mio/repository"
	"mio/service"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (TopicController) List(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	list, total, err := service.DefaultTopicService.GetTopicDetailPageList(repository.GetTopicPageListBy{
		ID:         form.ID,
		TopicTagId: form.TopicTagId,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
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
func (TopicController) GetShareWeappQrCode(c *gin.Context) (gin.H, error) {
	form := GetWeappQrCodeFrom{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := util.GetAuthUser(c)

	buffers, contType, err := service.DefaultTopicService.GetShareWeappQrCode(user.ID, form.TopicId)
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

	like, err := service.TopicLikeService{}.ChangeLikeStatus(form.TopicId, user.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"status": like.Status,
	}, nil
}
