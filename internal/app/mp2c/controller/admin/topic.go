package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	communityPdr "mio/internal/pkg/queue/producer/community"
	"mio/internal/pkg/queue/types/message/communitymsg"
	communityModel "mio/internal/pkg/repository/community"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/limit"
	"strconv"
	"strings"
	"time"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (ctr TopicController) List(c *gin.Context) (gin.H, error) {
	form := TopicListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	cond := community.AdminTopicListParams{
		ID:            form.ID,
		Title:         form.Title,
		UserId:        form.UserId,
		UserName:      form.UserName,
		Status:        form.Status,
		IsTop:         form.IsTop,
		IsEssence:     form.IsEssence,
		IsPartners:    form.IsPartners,
		Position:      form.Position,
		PushStartTime: form.PushStartTime,
		PushEndTime:   form.PushEndTime,
		Type:          form.Type,
		ActivityType:  form.ActivityType,
		Offset:        form.Offset(),
		Limit:         form.Limit(),
	}

	tagIds := strings.Split(form.TagId, ",")

	if len(tagIds) == 1 {
		i64, _ := strconv.ParseInt(tagIds[0], 10, 64)
		cond.TagId = i64
	} else if len(tagIds) > 1 {
		cond.TagIds = tagIds
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := community.NewTopicAdminService(ctx)
	activitiesSignupService := community.NewCommunityActivitiesSignupService(ctx)
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
	countList, err := activitiesSignupService.FindListCount(community.FindListCountReq{TopicIds: ids})
	if err != nil {
		return nil, err
	}
	signupCountMap := make(map[int64]int64, 0)
	for _, item := range countList {
		signupCountMap[item.TopicId] = item.NumOfSignup
	}
	rootCommentCount := adminTopicService.GetCommentCount(ids)
	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	for _, item := range list {
		item.CommentCount = topic2comment[item.Id]
		item.Activity.NumOfSignup = int(signupCountMap[item.Id])
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
	adminTopicService := community.NewTopicAdminService(ctx)

	topic, err := adminTopicService.DetailTopic(form.ID)
	if err != nil {
		return nil, err
	}

	if topic.Type == communityModel.TopicTypeActivity {
		topic.Activity.Status = 1
		if topic.Activity.SignupDeadline.Before(time.Now()) {
			topic.Activity.Status = 2
		}
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
	adminTopicService := community.NewTopicAdminService(ctx)

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
	adminTopicService := community.NewTopicAdminService(ctx)

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
	adminTopicService := community.NewTopicAdminService(ctx)
	messageService := message.NewWebMessageService(ctx)

	topic, err := adminTopicService.SoftDeleteTopic(form.ID, form.Reason)
	if err != nil {
		return nil, err
	}
	//发消息
	err = messageService.SendMessage(message.SendWebMessage{
		SendId:       0,
		RecId:        topic.UserId,
		Key:          "down_topic",
		TurnId:       topic.Id,
		TurnType:     message.MsgTurnTypeArticle,
		Type:         message.MsgTypeSystem,
		MessageNotes: topic.Title,
	})

	if err != nil {
		app.Logger.Errorf("【帖子下架】站内信发送失败:%s", err.Error())
	}

	return nil, nil
}

func (ctr TopicController) Down(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	//更新帖子
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := community.NewTopicAdminService(ctx)
	topic, err := adminTopicService.DownTopic(form.ID, form.Reason)
	if err != nil {
		return nil, err
	}
	ctr.seekingStore("down", topic)
	ctr.sendMessage(ctx, "down_topic", 0, topic.UserId, topic.Id, topic.Title)
	return nil, nil
}

// Review 审核
func (ctr TopicController) Review(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminTopicService := community.NewTopicAdminService(ctx)

	topic, _, err := adminTopicService.Review(form.ID, form.Status, form.Reason)
	if err != nil {
		return nil, err
	}

	pointService := service.NewPointService(ctx)
	key := ""
	if topic.Status == 3 {
		ctr.seekingStore("push", topic)
		//审核通过发积分
		keyPrefix := fmt.Sprintf("%s:%s", "periodLimit:sendPoint:article:push", topic.CreatedAt.Time.Format("2006-01-02"))
		PeriodLimit := limit.NewPeriodLimit(int(time.Hour.Seconds()*24), 2, app.Redis, keyPrefix, limit.PeriodAlign())
		resNumber, err := PeriodLimit.TakeCtx(ctx.Context, topic.User.OpenId)
		if err != nil {
			return nil, err
		}
		key = "push_topic_v2"
		if resNumber == 1 || resNumber == 2 {
			_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
				OpenId:       topic.User.OpenId,
				Type:         entity.POINT_ARTICLE,
				BizId:        util.UUID(),
				ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_ARTICLE]),
				AdminId:      0,
				Note:         fmt.Sprintf("审核帖子:%s；审核状态:%d", topic.Title, topic.Status),
				AdditionInfo: strconv.FormatInt(topic.Id, 10),
			})
			if err != nil {
				app.Logger.Errorf("【帖子审核】积分发放失败:%s; 帖子ID:%v; 用户ID:%v", err.Error(), topic.Id, topic.UserId)
			}
			key = "push_topic"
		}
	}

	if topic.Status == 4 {
		key = "down_topic"
	}
	if topic.Status == 2 {
		key = "fail_topic"
	}
	if key != "" {
		ctr.sendMessage(ctx, key, 0, topic.UserId, topic.Id, "")
		ctr.zhuGe(int(topic.Status), topic.Type, topic.User.OpenId)
		ctr.sensors(int(topic.Status), topic.Type, topic.User.OpenId, topic.Id)
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
	adminTopicService := community.NewTopicAdminService(ctx)
	messageService := message.NewWebMessageService(ctx)

	topic, isFirst, err := adminTopicService.Top(form.ID, form.IsTop)

	if err != nil {
		return nil, err
	}

	if isFirst {
		//发消息
		err = messageService.SendMessage(message.SendWebMessage{
			SendId:   0,
			RecId:    topic.UserId,
			Key:      "top_topic_v2",
			Type:     message.MsgTypeSystem,
			TurnType: message.MsgTurnTypeArticle,
			TurnId:   topic.Id,
		})

		if err != nil {
			app.Logger.Errorf("【帖子置顶】站内信发送失败:%s", err.Error())
		}

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
	adminTopicService := community.NewTopicAdminService(ctx)
	pointService := service.NewPointService(ctx)
	messageService := message.NewWebMessageService(ctx)

	topic, isFirst, err := adminTopicService.Essence(form.ID, form.IsEssence)
	if err != nil {
		return nil, err
	}

	if topic.IsEssence == 1 {
		keyPrefix := "periodLimit:sendPoint:article:essence:"
		PeriodLimit := limit.NewPeriodLimit(int(time.Hour.Seconds()*24), 2, app.Redis, keyPrefix, limit.PeriodAlign())
		resNumber, err := PeriodLimit.TakeCtx(ctx.Context, topic.User.OpenId)

		if err != nil {
			return nil, err
		}

		key := "essence_topic_v2"
		if resNumber == 1 || resNumber == 2 {
			_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
				OpenId:       topic.User.OpenId,
				Type:         entity.POINT_RECOMMEND,
				BizId:        util.UUID(),
				ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_RECOMMEND]),
				AdminId:      0,
				Note:         "帖子 \"" + topic.Title + "...\" 被设为精华",
				AdditionInfo: strconv.FormatInt(topic.Id, 10),
			})
			if err != nil {
				app.Logger.Errorf("【帖子加精】积分增加失败:%s; 帖子ID:%v; 用户ID:%v", err.Error(), topic.Id, topic.UserId)
			}
			key = "essence_topic"
		}

		//发消息
		if isFirst {
			err = messageService.SendMessage(message.SendWebMessage{
				SendId:   0,
				RecId:    topic.UserId,
				Key:      key,
				Type:     message.MsgTypeSystem,
				TurnType: message.MsgTurnTypeArticle,
				TurnId:   topic.Id,
			})

			if err != nil {
				app.Logger.Errorf("【精华帖子】站内信发送失败:%s", err.Error())
			}
		}
	}

	return nil, nil
}

func (ctr TopicController) seekingStore(event string, topic *entity.Topic) {
	err := communityPdr.SeekingStore(communitymsg.Topic{
		Event:     event,
		Id:        topic.Id,
		UserId:    topic.UserId,
		Status:    int(topic.Status),
		Type:      topic.Type,
		Tags:      topic.Tags,
		CreatedAt: topic.CreatedAt.Time,
	})
	if err != nil {
		app.Logger.Errorf("[城市碳秘] communityPdr Err: %s", err.Error())
	}
}

func (ctr TopicController) sendMessage(ctx *context.MioContext, key string, sendId, recId, turnId int64, msg string) {
	messageService := message.NewWebMessageService(ctx)
	err := messageService.SendMessage(message.SendWebMessage{
		SendId:       sendId,
		RecId:        recId,
		Key:          key,
		Type:         message.MsgTypeSystem,
		TurnType:     message.MsgTurnTypeArticle,
		TurnId:       turnId,
		MessageNotes: msg,
	})
	if err != nil {
		app.Logger.Errorf("【帖子审核】站内信发送失败:%s", err.Error())
	}
}

func (ctr TopicController) zhuGe(status, tp int, openId string) {
	var review string
	switch status {
	case 3:
		review = "审核通过"
	default:
		review = "审核未通过"
	}
	//诸葛打点
	zhuGeAttr := make(map[string]interface{}, 0)
	zhuGeAttr["场景"] = "发布帖子"
	zhuGeAttr["审核状态"] = review
	eventName := config.ZhuGeEventName.PostArticle
	if tp == 1 {
		zhuGeAttr["场景"] = "发布活动"
		eventName = config.ZhuGeEventName.PostActivity
	}
	track.DefaultZhuGeService().Track(eventName, openId, zhuGeAttr)
}

func (ctr TopicController) sensors(status, tp int, openId string, topicId int64) {
	var review string
	switch status {
	case 3:
		review = "审核通过"
	default:
		review = "审核未通过"
	}
	//诸葛打点
	scene := "发布帖子"
	if tp == 1 {
		scene = "发布活动"
	}
	track.DefaultSensorsService().Track(true, config.SensorsEventName.CommunityTopic, openId, map[string]interface{}{
		"status":   review,
		"type":     scene,
		"topic_id": int(topicId),
	})
}
