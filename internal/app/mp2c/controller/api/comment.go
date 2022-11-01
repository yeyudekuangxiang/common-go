package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"strconv"
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
	commentService := service.NewCommentService(ctx)
	commentLikeService := service.NewCommentLikeService(ctx)

	list, total, err := commentService.FindListAndChild(&req, form.Offset(), form.Limit())
	if err != nil {
		return nil, err
	}
	//获取点赞记录
	commentRes := make([]*entity.CommentRes, 0)
	likeMap := make(map[int64]int, 0)
	commentLike := commentLikeService.GetLikeInfoByUser(user.ID)
	if len(commentLike) > 0 {
		for _, item := range commentLike {
			likeMap[item.CommentId] = int(item.Status)
		}
	}

	for _, item := range list {
		res := item.CommentRes()
		if _, ok := likeMap[item.Id]; ok {
			res.IsLike = likeMap[item.Id]
		}
		if item.RootChild != nil {
			for _, childItem := range item.RootChild {
				childRes := childItem.CommentRes()
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
	commentLikeService := service.NewCommentLikeService(ctx)
	commentService := service.NewCommentService(ctx)

	data := &entity.CommentIndex{
		RootCommentId: form.ID,
	}

	list, total, err := commentService.FindSubList(data, form.Offset(), form.Limit())
	if err != nil {
		return nil, err
	}
	//获取点赞记录
	commentRes := make([]*entity.CommentRes, 0)
	likeMap := make(map[int64]int, 0)
	commentLike := commentLikeService.GetLikeInfoByUser(user.ID)
	if len(commentLike) > 0 {
		for _, item := range commentLike {
			likeMap[item.CommentId] = int(item.Status)
		}
	}

	for _, item := range list {
		res := item.CommentRes()
		if _, ok := likeMap[item.Id]; ok {
			res.IsLike = likeMap[item.Id]
		}
		if item.RootChild != nil {
			for _, childItem := range item.RootChild {
				childRes := childItem.CommentRes()
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

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := service.NewCommentService(ctx)
	messageService := message.NewWebMessageService(ctx)

	comment, recId, err := commentService.CreateComment(user.ID, form.ObjId, form.Root, form.Parent, form.Message, user.OpenId)
	if err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}

	//更新积分
	msg := comment.Message
	messagerune := []rune(comment.Message)
	if len(messagerune) > 8 {
		msg = string(messagerune[0:8])
	}

	point := int64(entity.PointCollectValueMap[entity.POINT_COMMENT])
	pointService := service.NewPointService(ctx)
	_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       user.OpenId,
		Type:         entity.POINT_COMMENT,
		BizId:        util.UUID(),
		ChangePoint:  point,
		AdminId:      0,
		Note:         "评论" + msg + "..." + "成功",
		AdditionInfo: strconv.FormatInt(form.ObjId, 10) + "#" + strconv.FormatInt(comment.Id, 10),
	})
	if err != nil {
		point = 0
	}

	//发送消息
	msgKey := "reply_topic"
	recObjId := form.ObjId
	if form.Parent != 0 {
		msgKey = "reply_comment"
		recObjId = form.Parent
	}

	err = messageService.SendMessage(message.SendWebMessage{
		SendId:   user.ID,
		RecId:    recId,
		Key:      msgKey,
		RecObjId: recObjId,
		Type:     1,
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
	user := apiutil.GetAuthUser(c)
	if user.Auth == 0 {
		return gin.H{"comment": nil, "point": 0}, errno.ErrCommon.WithMessage("无权限")
	}
	form := CommentEditForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}
	if form.Message != "" {
		//检查内容
		if err := validator.CheckMsgWithOpenId(user.OpenId, form.Message); err != nil {
			app.Logger.Error(fmt.Errorf("update Comment error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "更新评论"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
			return gin.H{"comment": nil, "point": 0}, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := service.NewCommentService(ctx)

	err := commentService.UpdateComment(user.ID, form.CommentId, form.Message)
	if err != nil {
		return gin.H{}, err
	}
	return gin.H{}, nil
}

func (ctr *CommentController) Delete(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := service.NewCommentService(ctx)

	err := commentService.DelCommentSoft(user.ID, form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (ctr *CommentController) Detail(c *gin.Context) (gin.H, error) {
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, nil
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	commentService := service.NewCommentService(ctx)

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
	commentService := service.NewCommentService(ctx)
	messageService := message.NewWebMessageService(ctx)

	like, point, recId, err := commentService.Like(user.ID, form.CommentId, user.OpenId)
	if err != nil {
		return nil, err
	}

	err = messageService.SendMessage(message.SendWebMessage{
		SendId:   user.ID,
		RecId:    recId,
		Key:      "like_comment",
		RecObjId: form.CommentId,
		Type:     1,
	})

	if err != nil {
		app.Logger.Errorf("【评论点赞】站内信发送失败:%s", err.Error())
	}

	return gin.H{
		"status": like.Status,
		"point":  point,
	}, nil
}
