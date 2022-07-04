package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCommentController = CommentController{}

type CommentController struct {
}

// List 除去置顶外，按时间排序
func (CommentController) List(c *gin.Context) (gin.H, error) {
	form := CommentForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	req := entity.CommentIndex{
		ObjId: form.TopicId,
	}
	list, total, err := service.DefaultCommentService.FindPageListByPage(&req, form.Page, form.PageSize, "like_count DESC")
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

// SubCommentList 根据root分页获取子评论
func (CommentController) SubCommentList(c *gin.Context) (gin.H, error) {
	form := SubCommentForm{}
	if err := apiutil.BindForm(c, form); err != nil {
		return nil, err
	}
	//data := entity.CommentIndex{
	//	Parent: form.CommentId,
	//}
	//list, total, err := service.DefaultCommentService.FindPageListByPage(&data)
	return gin.H{}, nil
}

func (CommentController) Create(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := CommentCreateForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}
	err := service.DefaultCommentService.CreateComment(user.ID, form.ObjId, form.Root, form.Parent, form.Message)
	if err != nil {
		return gin.H{}, err
	}
	return gin.H{}, nil
}

func (CommentController) Update(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := CommentEditForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}
	err := service.DefaultCommentService.UpdateComment(user.ID, form.CommentId, form.Message)
	if err != nil {
		return gin.H{}, err
	}
	return gin.H{}, nil
}

func (CommentController) Delete(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, err
	}
	err := service.DefaultCommentService.DelCommentSoft(user.ID, form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (CommentController) Detail(c *gin.Context) (gin.H, error) {
	//user := apiutil.GetAuthUser(c)
	form := IdForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{}, nil
	}
	one, err := service.DefaultCommentService.FindOne(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"result": one,
	}, nil
}
