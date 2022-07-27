package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

//List 获取文章列表
func (ctr *TopicController) List(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	list, total, err := service.DefaultTopicService.GetTopicDetailPageList(repository.GetTopicPageListBy{
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

func (ctr *TopicController) ListFlow(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	list, total, err := service.DefaultTopicService.GetTopicDetailPageListByFlow(repository.GetTopicPageListBy{
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

//GetShareWeappQrCode 获取分享二维码
func (ctr *TopicController) GetShareWeappQrCode(c *gin.Context) (gin.H, error) {
	form := GetWeappQrCodeFrom{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	buffers, contType, err := service.DefaultTopicService.GetShareWeappQrCode(int(user.ID), form.TopicId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"buffers":     buffers,
		"contentType": contType,
	}, nil
}

//ChangeTopicLike 点赞 / 取消点赞
func (ctr *TopicController) ChangeTopicLike(c *gin.Context) (gin.H, error) {
	form := ChangeTopicLikeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	like, err := service.TopicLikeService{}.ChangeLikeStatus(form.TopicId, int(user.ID))
	if err != nil {
		return nil, err
	}
	return gin.H{
		"status": like.Status,
	}, nil
}

func (ctr *TopicController) CreateTopic(c *gin.Context) (gin.H, error) {
	form := CreateTopicForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	//创建帖子
	err := service.DefaultTopicService.CreateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *TopicController) UpdateTopic(c *gin.Context) (gin.H, error) {
	form := UpdateTopicForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user
	user := apiutil.GetAuthUser(c)
	//更新帖子
	err := service.DefaultTopicService.UpdateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.ID, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *TopicController) DelTopic(c *gin.Context) (gin.H, error) {
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user
	user := apiutil.GetAuthUser(c)
	//更新帖子
	err := service.DefaultTopicService.DelTopic(user.ID, form.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *TopicController) DetailTopic(c *gin.Context) (gin.H, error) {
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//获取帖子
	topic, err := service.DefaultTopicService.DetailTopic(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr *TopicController) MyTopic(c *gin.Context) (gin.H, error) {
	return nil, nil
}

func (ctr *TopicController) MyTopicList(context *gin.Context) (gin.H, error) {
	return nil, nil
}
