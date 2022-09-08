package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller"
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

	page := "pages/community/details/index"
	scene := fmt.Sprintf("tid=%d&uid=%d&s=p&model=coolmio", form.TopicId, user.ID)

	qr, err := service.NewQRCodeService().GetUnlimitedQRCodeRaw(page, scene, 100)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"qrcode": qr,
	}, nil
}

//ChangeTopicLike 点赞 / 取消点赞
func (ctr *TopicController) ChangeTopicLike(c *gin.Context) (gin.H, error) {
	form := ChangeTopicLikeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	like, point, err := service.TopicLikeService{}.ChangeLikeStatus(form.TopicId, int(user.ID), user.OpenId)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"point":  point,
		"status": like.Status,
	}, nil
}

//ListTopic 帖子列表+顶级评论+顶级评论下子评论3条
func (ctr *TopicController) ListTopic(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	list, total, err := service.DefaultTopicService.GetTopicList(repository.GetTopicPageListBy{
		TopicTagId: form.TopicTagId,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	resList := make([]*entity.TopicItemRes, 0)
	//点赞数据
	likeMap := make(map[int]int, 0)
	likeList, _ := service.DefaultTopicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = int(item.Status)
		}
	}

	//获取顶级评论数量
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	rootCommentCount := service.DefaultTopicService.GetRootCommentCount(ids)
	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	for _, item := range list {
		res := item.TopicItemRes()
		res.CommentCount = topic2comment[res.Id]
		if _, ok := likeMap[int(res.Id)]; ok {
			res.IsLike = likeMap[int(res.Id)]
		}
		resList = append(resList, res)
	}
	app.Logger.Infof("GetTopicDetailPageListByFlow user:%d form:%+v ids:%+v", user.ID, form, ids)
	return gin.H{
		"list":     resList,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

//CreateTopic 创建帖子
func (ctr *TopicController) CreateTopic(c *gin.Context) (gin.H, error) {
	form := CreateTopicForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errors.New("无权限")
	}
	//创建帖子
	topic, err := service.DefaultTopicService.CreateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
		"point": 0,
	}, nil
}

func (ctr *TopicController) UpdateTopic(c *gin.Context) (gin.H, error) {
	form := UpdateTopicForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errors.New("无权限")
	}
	//更新帖子
	topic, err := service.DefaultTopicService.UpdateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.ID, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr *TopicController) DelTopic(c *gin.Context) (gin.H, error) {
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errors.New("无权限")
	}
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
	user := apiutil.GetAuthUser(c)
	//获取帖子
	topic, err := service.DefaultTopicService.DetailTopic(form.ID)
	if err != nil {
		return nil, err
	}
	topicRes := topic.TopicItemRes()
	//获取评论数量

	CommentCount := service.DefaultTopicService.GetCommentCount([]int64{topicRes.Id})
	//组装数据---帖子的顶级评论数量
	if len(CommentCount) > 0 {
		topicRes.CommentCount = CommentCount[0].Total
	}
	//获取点赞数据
	like, err := service.DefaultTopicLikeService.GetOneByTopic(topic.Id, user.ID)
	if err == nil {
		topicRes.IsLike = int(like.Status)
	}
	return gin.H{
		"topic": topicRes,
	}, nil
}

func (ctr *TopicController) MyTopic(c *gin.Context) (gin.H, error) {
	form := controller.PageFrom{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)
	list, total, err := service.DefaultTopicService.GetMyTopicList(repository.GetTopicPageListBy{
		UserId: user.ID,
		Limit:  form.Limit(),
		Offset: form.Offset(),
	})
	if err != nil {
		return nil, err
	}
	resList := make([]*entity.TopicItemRes, 0)
	//点赞数据
	likeMap := make(map[int]int, 0)
	likeList, _ := service.DefaultTopicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = int(item.Status)
		}
	}

	//获取顶级评论数量
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	rootCommentCount := service.DefaultTopicService.GetRootCommentCount(ids)
	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	for _, item := range list {
		res := item.TopicItemRes()
		res.CommentCount = topic2comment[res.Id]
		if _, ok := likeMap[int(res.Id)]; ok {
			res.IsLike = likeMap[int(res.Id)]
		}
		resList = append(resList, res)
	}
	return gin.H{
		"list":     resList,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, err
}
