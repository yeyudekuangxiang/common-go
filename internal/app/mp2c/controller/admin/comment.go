package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
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
	list, total, err := service.DefaultCommentAdminService.CommentList(form.Comment, form.UserId, form.TopicId, form.Limit(), form.Offset())
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
	if err := service.DefaultCommentAdminService.DelCommentSoft(form.ID, form.Reason); err != nil {
		return nil, err
	}
	return nil, nil
}
