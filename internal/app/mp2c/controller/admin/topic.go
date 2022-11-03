package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"strconv"
	"strings"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (ctr TopicController) List(c *gin.Context) (gin.H, error) {
	form := TopicListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	cond := repository.TopicListRequest{
		ID:         form.ID,
		Title:      form.Title,
		UserId:     form.UserId,
		UserName:   form.UserName,
		Status:     form.Status,
		IsTop:      form.IsTop,
		IsEssence:  form.IsEssence,
		IsPartners: form.IsPartners,
		Position:   form.Position,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
	}

	tagIds := strings.Split(form.TagId, ",")

	if len(tagIds) == 1 {
		float, _ := strconv.ParseInt(tagIds[0], 10, 64)
		cond.TagId = float
	} else if len(tagIds) > 1 {
		cond.TagIds = tagIds
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)

	//get topic by params
	list, total, err := adminTopicService.GetTopicList(cond)

	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	//获取顶级评论数量
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	rootCommentCount := adminTopicService.GetCommentCount(ids)
	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	for _, item := range list {
		item.CommentCount = topic2comment[item.Id]
	}
	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (ctr TopicController) Detail(c *gin.Context) (gin.H, error) {
	form := IDForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)

	topic, err := adminTopicService.DetailTopic(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr TopicController) Create(c *gin.Context) (gin.H, error) {
	form := CreateTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//创建帖子
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)

	err := adminTopicService.CreateTopic(int64(1451), form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

func (ctr TopicController) Update(c *gin.Context) (gin.H, error) {
	form := UpdateTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//更新帖子
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)

	err := adminTopicService.UpdateTopic(form.ID, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Delete 软删除
func (ctr TopicController) Delete(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	//更新帖子
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)

	if err := adminTopicService.DeleteTopic(form.ID, form.Reason); err != nil {
		return nil, err
	}
	return nil, nil
}

// Review 审核
func (ctr TopicController) Review(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)

	topic, err := adminTopicService.Review(form.ID, form.Status, form.Reason)
	if err != nil {
		return nil, err
	}

	pointService := service.NewPointService(ctx)
	pointTransService := service.NewPointTransactionService(ctx)
	messageService := message.NewWebMessageService(ctx)

	user, _ := service.DefaultUserService.GetUserById(topic.UserId)

	by, err := pointTransService.FindBy(repository.FindPointTransactionBy{
		OpenId:       user.OpenId,
		Type:         string(entity.POINT_ARTICLE),
		AdditionInfo: strconv.FormatInt(topic.Id, 10),
	})

	var point int64
	if topic.Status == 3 && form.Status == 4 {
		point = -int64(entity.PointCollectValueMap[entity.POINT_ARTICLE])
	}

	title := topic.Title
	if len([]rune(title)) > 8 {
		title = string([]rune(title)[0:8]) + "..."
	}

	if form.Status == 3 && by.ID == 0 {
		point = int64(entity.PointCollectValueMap[entity.POINT_ARTICLE])
		_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_ARTICLE,
			BizId:        util.UUID(),
			ChangePoint:  point,
			AdminId:      0,
			Note:         fmt.Sprintf("审核笔记:%s；审核状态:%d", topic.Title, topic.Status),
			AdditionInfo: strconv.FormatInt(topic.Id, 10),
		})
	}

	//发消息
	var key string
	var t int
	if form.Status == 3 {
		key = "post_topic"
		t = 4
	}
	if form.Status == 4 {
		key = "fail_topic"
		t = 6
	}
	err = messageService.SendMessage(message.SendWebMessage{
		SendId:   0,
		RecId:    user.ID,
		Key:      key,
		TurnId:   topic.Id,
		TurnType: 1,
		Type:     t,
	})

	if err != nil {
		app.Logger.Errorf("【文章审核】站内信发送失败:%s", err.Error())
	}
	return nil, nil
}

// Top 置顶
func (ctr TopicController) Top(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)
	messageService := message.NewWebMessageService(ctx)

	topic, err := adminTopicService.Top(form.ID, form.IsTop)
	if err != nil {
		return nil, err
	}

	//发消息
	err = messageService.SendMessage(message.SendWebMessage{
		SendId:   0,
		RecId:    topic.UserId,
		Key:      "top_topic",
		TurnType: 1,
		TurnId:   topic.Id,
		Type:     5,
	})

	if err != nil {
		app.Logger.Errorf("【文章置顶】站内信发送失败:%s", err.Error())
	}

	return nil, nil
}

// Essence 精华
func (ctr TopicController) Essence(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := service.NewTopicAdminService(ctx)
	messageService := message.NewWebMessageService(ctx)

	topic, err := adminTopicService.Essence(form.ID, form.IsEssence)
	if err != nil {
		return nil, err
	}

	//发放积分
	//发放积分
	if form.IsEssence == 1 {
		title := topic.Title
		if len([]rune(title)) > 8 {
			title = string([]rune(title)[0:8]) + "..."
		}
		user, _ := service.DefaultUserService.GetUserById(topic.UserId)
		pointService := service.NewPointService(context.NewMioContext())
		_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_RECOMMEND,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_RECOMMEND]),
			AdminId:      0,
			Note:         "笔记 \"" + title + "...\" 被设为精华",
			AdditionInfo: strconv.FormatInt(topic.Id, 10),
		})
	}

	//发消息
	err = messageService.SendMessage(message.SendWebMessage{
		SendId:   0,
		RecId:    topic.UserId,
		Key:      "essence_topic",
		TurnType: 1,
		TurnId:   topic.Id,
		Type:     5,
	})

	if err != nil {
		app.Logger.Errorf("【精华文章】站内信发送失败:%s", err.Error())
	}

	return nil, nil
}
