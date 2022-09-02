package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
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
	list, total, err := service.DefaultCommentService.FindListAndChild(&req, form.Offset(), form.Limit())
	if err != nil {
		return nil, err
	}
	//获取点赞记录
	commentRes := make([]*entity.CommentRes, 0)
	likeMap := make(map[int64]int, 0)
	commentLike := service.DefaultCommentLikeService.GetLikeInfoByUser(user.ID)
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

	data := &entity.CommentIndex{
		RootCommentId: form.ID,
	}
	list, total, err := service.DefaultCommentService.FindSubList(data, form.Offset(), form.Limit())
	if err != nil {
		return nil, err
	}
	//获取点赞记录
	commentRes := make([]*entity.CommentRes, 0)
	likeMap := make(map[int64]int, 0)
	commentLike := service.DefaultCommentLikeService.GetLikeInfoByUser(user.ID)
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
	user := apiutil.GetAuthUser(c)
	form := CommentCreateForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}
	comment, point, err := service.DefaultCommentService.CreateComment(user.ID, form.ObjId, form.Root, form.Parent, form.Message, user.OpenId)
	if err != nil {
		return gin.H{"comment": nil, "point": 0}, err
	}
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

	like, err := service.DefaultCommentService.Like(user.ID, form.CommentId, user.OpenId)
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
