package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/kumiaoCommunity"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strconv"
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

	list, total, err := kumiaoCommunity.DefaultTopicService.GetTopicDetailPageList(repository.GetTopicPageListBy{
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

//func (ctr *TopicController) ListFlow(c *gin.Context) (gin.H, error) {
//	form := GetTopicPageListForm{}
//	if err := apiutil.BindForm(c, &form); err != nil {
//		return nil, err
//	}
//
//	user := apiutil.GetAuthUser(c)
//
//	list, total, err := kumiaoCommunity.DefaultTopicService.GetTopicDetailPageListByFlow(repository.GetTopicPageListBy{
//		ID:         form.ID,
//		TopicTagId: form.TopicTagId,
//		Offset:     form.Offset(),
//		Limit:      form.Limit(),
//		UserId:     user.ID,
//	})
//	if err != nil {
//		return nil, err
//	}
//	ids := make([]int64, 0)
//	for _, item := range list {
//		ids = append(ids, item.Id)
//	}
//	app.Logger.Infof("user:%d form:%+v ids:%+v", user.ID, form, ids)
//
//	return gin.H{
//		"list":     list,
//		"total":    total,
//		"page":     form.Page,
//		"pageSize": form.PageSize,
//	}, nil
//}

//GetShareWeappQrCode 获取分享二维码
func (ctr *TopicController) GetShareWeappQrCode(c *gin.Context) (gin.H, error) {
	form := GetWeappQrCodeFrom{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	page := "pages/community/details/index"
	scene := fmt.Sprintf("tid=%d&uid=%d&s=p&m=c", form.TopicId, user.ID)

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
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	topicLikeService := kumiaoCommunity.NewTopicLikeService(ctx)
	messageService := message.NewWebMessageService(ctx)

	resp, err := topicLikeService.ChangeLikeStatus(form.TopicId, user.ID, user.OpenId)
	if err != nil {
		return nil, err
	}

	title := resp.TopicTitle
	if len([]rune(title)) > 8 {
		title = string([]rune(title)[0:8]) + "..."
	}

	var point int64
	if resp.LikeStatus == 1 && resp.IsFirst == true {
		pointService := service.NewPointService(ctx)
		_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_LIKE,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_LIKE]),
			AdminId:      0,
			Note:         "为文章 \"" + title + "\" 点赞",
			AdditionInfo: strconv.FormatInt(resp.TopicId, 10),
		})

		if err == nil {
			point = int64(entity.PointCollectValueMap[entity.POINT_LIKE])
		}

		//发送消息
		err = messageService.SendMessage(message.SendWebMessage{
			SendId:   user.ID,
			RecId:    resp.TopicUserId,
			Key:      "like_topic",
			TurnType: 1,
			TurnId:   resp.TopicId,
			Type:     1,
			ShowId:   resp.TopicId,
		})
		if err != nil {
			app.Logger.Errorf("文章点赞站内信发送失败:%s", err.Error())
		}

	}

	return gin.H{
		"point":  point,
		"status": resp.LikeStatus,
	}, nil
}

//ListTopic 帖子列表+顶级评论+顶级评论下子评论3条
func (ctr *TopicController) ListTopic(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	list, total, err := kumiaoCommunity.DefaultTopicService.GetTopicList(repository.GetTopicPageListBy{
		TopicTagId: form.TopicTagId,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
		Order:      form.Order,
	})
	if err != nil {
		return nil, err
	}
	resList := make([]*entity.TopicItemRes, 0)
	//点赞数据
	likeMap := make(map[int64]struct{}, 0)
	topicLikeService := kumiaoCommunity.NewTopicLikeService(ctx)
	likeList, _ := topicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = struct{}{}
		}
	}
	//收藏数据
	collectionMap := make(map[int64]struct{}, 0)
	collectionService := kumiaoCommunity.NewCollectionService(ctx)
	collectionIds := collectionService.Collections(user.OpenId, 0, 0, 0)
	for _, collectionId := range collectionIds {
		collectionMap[collectionId] = struct{}{}
	}

	//评论数据
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}

	rootCommentCount := kumiaoCommunity.DefaultTopicService.GetRootCommentCount(ids)
	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	for _, item := range list {
		res := item.TopicItemRes()
		res.CommentCount = topic2comment[res.Id]
		if _, ok := likeMap[res.Id]; ok {
			res.IsLike = 1
		}
		if _, ok := collectionMap[res.Id]; ok {
			res.IsCollection = 1
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
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}
	form := CreateTopicForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//创建帖子
	topic, err := kumiaoCommunity.DefaultTopicService.CreateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
		"point": 0,
	}, nil
}

func (ctr *TopicController) UpdateTopic(c *gin.Context) (gin.H, error) {
	//user
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}

	form := UpdateTopicForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	//更新帖子
	topic, err := kumiaoCommunity.DefaultTopicService.UpdateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.ID, form.Title, form.Content, form.TagIds, form.Images)
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
		return nil, errno.ErrCommon.WithMessage("无权限")
	}
	//更新帖子
	err := kumiaoCommunity.DefaultTopicService.DelTopic(user.ID, form.ID)
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

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	topicService := kumiaoCommunity.NewTopicService(ctx)
	topicLikeService := kumiaoCommunity.NewTopicLikeService(ctx)
	collectService := kumiaoCommunity.NewCollectionService(ctx)

	//获取帖子
	topic, err := topicService.DetailTopic(form.ID)
	if err != nil {
		return nil, err
	}
	topicRes := topic.TopicItemRes()
	//获取评论数量

	CommentCount := topicService.GetCommentCount([]int64{topicRes.Id})
	// 组装数据
	// 评论
	if len(CommentCount) > 0 {
		topicRes.CommentCount = CommentCount[0].Total
	}

	// 点赞
	like, err := topicLikeService.GetOneByTopic(topic.Id, user.ID)
	if err == nil {
		topicRes.IsLike = int(like.Status)
	}
	// 收藏
	collection, err := collectService.FindOneByTopic(topic.Id, user.OpenId)

	if err == nil {
		topicRes.IsCollection = collection.Status
	}

	return gin.H{
		"topic": topicRes,
	}, nil
}

func (ctr *TopicController) MyTopic(c *gin.Context) (gin.H, error) {
	form := MyTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)
	status := form.Status
	if form.UserId != 0 {
		u, b, err := service.DefaultUserService.GetUserByID(form.UserId)
		if err != nil {
			return nil, errno.ErrCommon
		}
		if !b {
			return nil, errno.ErrUserNotFound
		}
		user = *u
		status = 3
	}

	list, total, err := kumiaoCommunity.DefaultTopicService.GetMyTopicList(repository.GetTopicPageListBy{
		UserId: user.ID,
		Status: status,
		Limit:  form.Limit(),
		Offset: form.Offset(),
	})

	if err != nil {
		return nil, err
	}

	resList := make([]*entity.TopicItemRes, 0)

	//点赞数据
	likeMap := make(map[int64]struct{}, 0)
	topicLikeService := kumiaoCommunity.NewTopicLikeService(ctx)
	likeList, _ := topicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = struct{}{}
		}
	}

	//评论数据
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}

	//收藏数据
	collectionMap := make(map[int64]struct{}, 0)
	collectionService := kumiaoCommunity.NewCollectionService(ctx)
	collectionIds := collectionService.Collections(user.OpenId, 0, 0, 0)
	for _, collectionId := range collectionIds {
		collectionMap[collectionId] = struct{}{}
	}

	rootCommentCount := kumiaoCommunity.DefaultTopicService.GetRootCommentCount(ids)
	// 组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	// 组装数据---点赞数据 收藏数据
	for _, item := range list {
		res := item.TopicItemRes()
		res.CommentCount = topic2comment[res.Id]
		if _, ok := likeMap[res.Id]; ok {
			res.IsLike = 1
		}
		if _, ok := collectionMap[res.Id]; ok {
			res.IsCollection = 1
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
