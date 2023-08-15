package community

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/queue/producer/growth_system"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/limit"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultCommentController = &CommentController{}

type CommentController struct {
}

// RootList 分页获取顶级评论 及 每条顶级评论下3条子评论
func (ctr *CommentController) RootList(c *gin.Context) (gin.H, error) {

	form := ListFormById{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	req := entity.CommentIndex{
		ObjId: form.ID,
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)
	commentLikeService := community.NewCommentLikeService(ctx)

	list, total, err := commentService.FindListAndChild(&req, form.Offset(), form.Limit())
	if err != nil {
		return nil, err
	}

	//获取点赞记录
	likeMap := make(map[int64]int, 0)
	commentLike := commentLikeService.GetLikeInfoByUser(user.ID)
	if len(commentLike) > 0 {
		for _, item := range commentLike {
			likeMap[item.CommentId] = int(item.Status)
		}
	}

	commentRes := make([]*entity.Comment, 0)
	for _, item := range list {
		res := item.CommentResp()
		if _, ok := likeMap[item.Id]; ok {
			res.IsLike = likeMap[item.Id]
		}
		if item.RootChild != nil {
			for _, childItem := range item.RootChild {
				childRes := childItem.CommentResp()
				if _, ok := likeMap[childItem.Id]; ok {
					childRes.IsLike = likeMap[childItem.Id]
				}
				res.RootChild = append(res.RootChild, childRes)
			}
		}
		commentRes = append(commentRes, res)
	}
	return gin.H{
		"list":     commentRes,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

// SubList 根据顶级评论分页获取子评论
func (ctr *CommentController) SubList(c *gin.Context) (gin.H, error) {
	form := ListFormById{}
	if err := apiutil.BindForm(c, form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentLikeService := community.NewCommentLikeService(ctx)
	commentService := community.NewCommentService(ctx)

	data := &entity.CommentIndex{
		RootCommentId: form.ID,
	}

	list, total, err := commentService.FindSubList(data, form.Offset(), form.Limit())
	if err != nil {
		return nil, err
	}
	//获取点赞记录
	commentRes := make([]*entity.Comment, 0)
	likeMap := make(map[int64]int, 0)
	commentLike := commentLikeService.GetLikeInfoByUser(user.ID)
	if len(commentLike) > 0 {
		for _, item := range commentLike {
			likeMap[item.CommentId] = int(item.Status)
		}
	}

	for _, item := range list {
		res := item.CommentResp()
		if _, ok := likeMap[item.Id]; ok {
			res.IsLike = likeMap[item.Id]
		}
		if item.RootChild != nil {
			for _, childItem := range item.RootChild {
				childRes := childItem.CommentResp()
				if _, ok := likeMap[childItem.Id]; ok {
					childRes.IsLike = likeMap[childItem.Id]
				}
				res.RootChild = append(res.RootChild, childRes)
			}
		}
		commentRes = append(commentRes, res)
	}
	return gin.H{
		"list":     commentRes,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (ctr *CommentController) Create(c *gin.Context) (gin.H, error) {
	form := CommentCreateForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}

	user := apiutil.GetAuthUser(c)

	if user.Auth == 0 {
		return gin.H{"comment": nil, "point": 0}, errno.ErrCommon.WithMessage("无权限")
	}

	userPlatform, exist, err := service.DefaultUserService.FindOneUserPlatformByGuid(c.Request.Context(), user.GUID, entity.UserPlatformWxMiniApp)

	if err != nil {
		return nil, err
	}

	if form.Message != "" && exist {
		//检查内容
		if err := validator.CheckMsgWithOpenId(userPlatform.Openid, form.Message); err != nil {
			app.Logger.Error(fmt.Errorf("create Comment error:%s", err.Error()))

			track.DefaultSensorsService().Track(false, config.SensorsEventName.MsgSecCheck, user.GUID, map[string]interface{}{
				"scene": "发布评论",
				"error": err.Error(),
			})
			return nil, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)
	messageService := message.NewWebMessageService(ctx)

	comment, toComment, recId, topic, err := commentService.CreateComment(user.ID, form.ObjId, form.Root, form.Parent, form.Message)
	if err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}
	//成长体系
	growth_system.GrowthSystemCommunityComment(growthsystemmsg.GrowthSystemParam{
		TaskSubType: string(entity.POINT_COMMENT),
		UserId:      strconv.FormatInt(user.ID, 10),
		TaskValue:   1,
	})

	//埋点
	trackData := map[string]interface{}{
		"topic_id":              int(form.ObjId),
		"topic_title":           topic.Title,
		"topic_user_id":         int(topic.User.ID),
		"topic_user_nickname":   topic.User.Nickname,
		"topic_tags":            topic.TopicTagId,
		"topic_type":            topic.Type,
		"comment_push_time":     comment.CreatedAt.Format(timetool.TimeFormat),
		"comment_user_id":       int(user.ID),
		"comment_user_position": string(user.Position),
		"comment_user_partner":  int(user.Partners),
	}
	track.DefaultSensorsService().Track(false, config.SensorsEventName.CommunityComment, user.GUID, trackData)
	//更新积分
	msg := comment.Message
	messagerune := []rune(comment.Message)
	if len(messagerune) > 8 {
		msg = string(messagerune[0:8])
	}

	keyPrefix := "periodLimit:sendPoint:comment:push:"
	periodLimit := limit.NewPeriodLimit(int(time.Hour.Seconds()*24), 3, app.Redis, keyPrefix, limit.PeriodAlign())
	resNumber, err := periodLimit.TakeCtx(ctx.Context, user.OpenId)
	if err != nil {
		return nil, err
	}

	point := 0
	if resNumber == 1 || resNumber == 2 {
		pointService := service.NewPointService(ctx)
		_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_COMMENT,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_COMMENT]),
			AdminId:      0,
			Note:         "评论" + msg + "..." + "成功",
			AdditionInfo: strconv.FormatInt(form.ObjId, 10) + "#" + strconv.FormatInt(comment.Id, 10),
		})
	}

	//发送消息
	var notes, msgKey string

	if form.Parent == 0 {
		msgKey = "reply_topic"
		notes = topic.Title
	} else {
		msgKey = "reply_comment"
		toMsg := toComment.Message
		messagerune = []rune(comment.Message)
		if len(messagerune) > 8 {
			toMsg = string(messagerune[0:8])
		}
		notes = toMsg
	}

	err = messageService.SendMessage(message.SendWebMessage{
		SendId:       user.ID,
		RecId:        recId,
		Key:          msgKey,
		TurnType:     message.MsgTurnTypeArticleComment,
		TurnId:       comment.Id,
		Type:         message.MsgTypeComment,
		MessageNotes: notes,
	})

	if err != nil {
		app.Logger.Errorf("评论站内信发送失败:%s", err.Error())
	}

	return gin.H{
		"comment": comment,
		"point":   point,
	}, nil
}

func (ctr *CommentController) Update(c *gin.Context) (gin.H, error) {
	form := CommentEditForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}

	user := apiutil.GetAuthUser(c)
	if user.Auth == 0 {
		return gin.H{"comment": nil, "point": 0}, errno.ErrCommon.WithMessage("无权限")
	}

	userPlatform, exist, err := service.DefaultUserService.FindOneUserPlatformByGuid(c.Request.Context(), user.GUID, entity.UserPlatformWxMiniApp)
	if err != nil {
		return nil, err
	}

	if form.Message != "" && exist {
		//检查内容
		if err := validator.CheckMsgWithOpenId(userPlatform.Openid, form.Message); err != nil {
			app.Logger.Error(fmt.Errorf("update Comment error:%s", err.Error()))
			/*zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "更新评论"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.GUID, zhuGeAttr)
			*/
			track.DefaultSensorsService().Track(false, config.SensorsEventName.MsgSecCheck, user.GUID, map[string]interface{}{
				"scene": "更新评论",
				"error": err.Error(),
			})

			return gin.H{"comment": nil, "point": 0}, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)

	err = commentService.UpdateComment(user.ID, form.CommentId, form.Message)
	if err != nil {
		return gin.H{}, err
	}
	return gin.H{}, nil
}

func (ctr *CommentController) Delete(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)

	err := commentService.DelCommentSoft(user.ID, form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (ctr *CommentController) Detail(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, nil
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)

	one, err := commentService.FindOne(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"result": one,
	}, nil
}

func (ctr *CommentController) Like(c *gin.Context) (gin.H, error) {
	form := ChangeCommentLikeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)
	messageService := message.NewWebMessageService(ctx)

	resp, err := commentService.Like(user.ID, form.CommentId, user.OpenId)
	if err != nil {
		return nil, err
	}

	var point int64
	if resp.LikeStatus == 1 {
		if resp.IsFirst {
			pointService := service.NewPointService(context.NewMioContext())
			_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
				OpenId:       user.OpenId,
				Type:         entity.POINT_LIKE,
				BizId:        util.UUID(),
				ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_LIKE]),
				AdminId:      0,
				Note:         "评论 \"" + resp.CommentMessage + "\" 点赞",
				AdditionInfo: "commendId: " + strconv.FormatInt(resp.CommentId, 10),
			})

			if err == nil {
				point = int64(entity.PointCollectValueMap[entity.POINT_LIKE])
			}

			err = messageService.SendMessage(message.SendWebMessage{
				SendId:       user.ID,
				RecId:        resp.CommentUserId,
				Key:          "like_comment",
				Type:         message.MsgTypeLike,
				TurnType:     message.MsgTurnTypeArticleComment,
				TurnId:       resp.CommentId,
				MessageNotes: resp.CommentMessage,
			})

			if err != nil {
				app.Logger.Errorf("【评论点赞】站内信发送失败:%s", err.Error())
			}
		}
		//成长体系
		growth_system.GrowthSystemCommunityLike(growthsystemmsg.GrowthSystemParam{
			TaskSubType: string(entity.POINT_LIKE),
			UserId:      strconv.FormatInt(user.ID, 10),
			TaskValue:   1,
		})
	}

	return gin.H{
		"status": resp.LikeStatus,
		"point":  point,
	}, nil
}

func (ctr *CommentController) TurnComment(c *gin.Context) (gin.H, error) {
	form := TurnCommentRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := community.NewCommentService(ctx)

	comment, err := commentService.TurnComment(community.TurnCommentReq{
		UserId:   user.ID,
		TurnType: form.TurnType,
		TurnId:   form.TurnId,
	})

	if err != nil {
		return nil, err
	}

	return gin.H{
		"comment": comment,
	}, nil
}
