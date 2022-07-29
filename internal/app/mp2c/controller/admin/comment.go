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

	return nil, nil
}

func (ctr *CommentController) Delete(c *gin.Context) (gin.H, error) {
	form := IDForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	if err := service.DefaultCommentAdminService.DelCommentSoft(form.ID); err != nil {
		return nil, err
	}
	return nil, nil
}
