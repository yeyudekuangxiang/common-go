package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
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
	req := entity.CommentIndex{
		ObjId: form.ID,
	}
	list, total, err := service.DefaultCommentService.FindListAndChild(&req, form.Offset(), form.Limit())
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

// SubList 根据顶级评论分页获取子评论
func (ctr *CommentController) SubList(c *gin.Context) (gin.H, error) {
	form := ListFormById{}
	if err := apiutil.BindForm(c, form); err != nil {
		return nil, err
	}
	data := &entity.CommentIndex{
		RootCommentId: form.ID,
	}
	list, total, err := service.DefaultCommentService.FindSubList(data, form.Offset(), form.Limit())
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

func (ctr *CommentController) Create(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := CommentCreateForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}
	comment, err := service.DefaultCommentService.CreateComment(user.ID, form.ObjId, form.Root, form.Parent, form.Message)
	if err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}
	//发放积分
	point := int64(entity.PointCollectValueMap[entity.POINT_COMMENT])
	pointService := service.NewPointService(context.NewMioContext())
	_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       user.OpenId,
		Type:         entity.POINT_COMMENT,
		BizId:        util.UUID(),
		ChangePoint:  point,
		AdminId:      0,
		Note:         strconv.FormatInt(comment.ID, 10),
		AdditionInfo: strconv.FormatInt(form.ObjId, 10) + "#" + strconv.FormatInt(comment.ID, 10),
	})
	return gin.H{
		"comment": comment,
		"point":   point,
	}, nil
}

func (ctr *CommentController) Update(c *gin.Context) (gin.H, error) {
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

func (ctr *CommentController) Delete(c *gin.Context) (gin.H, error) {
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

func (ctr *CommentController) Detail(c *gin.Context) (gin.H, error) {
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

func (ctr *CommentController) Like(c *gin.Context) (gin.H, error) {
	form := ChangeCommentLikeForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	like, err := service.DefaultCommentService.Like(user.ID, form.CommentId)
	if err != nil {
		return nil, err
	}
	var point int64
	if like.Status == 1 {
		point = int64(entity.PointCollectValueMap[entity.POINT_LIKE])
	}
	return gin.H{
		"status": like.Status,
		"point":  point,
	}, nil
}
