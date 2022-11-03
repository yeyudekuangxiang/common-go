package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/util/apiutil"
	"time"
)

var DefaultCommentController = &CommentController{}

type CommentController struct {
	Date time.Time
}

func (ctr *CommentController) List(c *gin.Context) (gin.H, error) {
	form := CommentListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminCommentService := service.NewCommentAdminService(ctx)

	list, total, err := adminCommentService.CommentList(form.Comment, form.UserId, form.TopicId, form.Limit(), form.Offset())
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Limit(),
		"pageSize": form.Offset(),
	}, nil
}

func (ctr *CommentController) Delete(c *gin.Context) (gin.H, error) {
	form := CommentDeleteRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	adminCommentService := service.NewCommentAdminService(ctx)
	messageService := message.NewWebMessageService(ctx)

	comment, err := adminCommentService.DelCommentSoft(form.ID, form.Reason)
	if err != nil {
		return nil, err
	}

	//发消息
	err = messageService.SendMessage(message.SendWebMessage{
		SendId:   0,
		RecId:    comment.MemberId,
		Type:     6,
		Key:      "fail_comment",
		TurnId:   comment.Id,
		TurnType: 2,
	})

	if err != nil {
		app.Logger.Errorf("【审核评论】站内信发送失败:%s", err.Error())
	}

	return nil, nil
}
