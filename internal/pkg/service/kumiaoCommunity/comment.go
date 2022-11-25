package kumiaoCommunity

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"time"
)

type (
	CommentService interface {
		FindOne(commentId int64) (*entity.CommentIndex, error)
		FindOneQuery(data *entity.CommentIndex) (*entity.CommentIndex, error)
		FindAll(data *entity.CommentIndex) ([]*entity.CommentIndex, int64, error)
		FindSubList(data *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error)
		FindListAndChild(data *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error)
		CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string, openId string) (entity.CommentIndex, *entity.CommentIndex, int64, error)
		UpdateComment(userId, commentId int64, message string) error
		DelComment(userId, commentId int64) error
		DelCommentSoft(userId, commentId int64) error
		Like(userId, commentId int64, openId string) (CommentChangeLikeResp, error)
		AddCommentLikeCount(commentId int64, num int) error
		TurnComment(params TurnCommentReq) (*APICommentResp, error)
		//TurnObj(params TurnCommentReq) (*APICommentResp, int64, error)
	}
)

type defaultCommentService struct {
	ctx                            *context.MioContext
	commentModel                   repository.CommentModel
	commentLikeModel               repository.CommentLikeModel
	topicModel                     repository.TopicModel
	carbonCommentModel             repository.CarbonCommentModel
	carbonCommentLikeModel         repository.CarbonCommentLikeModel
	carbonSecondHandCommodityModel repository.CarbonSecondHandCommodityModel
}

func NewCommentService(ctx *context.MioContext) CommentService {
	return &defaultCommentService{
		ctx:                            ctx,
		commentModel:                   repository.NewCommentModel(ctx),
		commentLikeModel:               repository.NewCommentLikeRepository(ctx),
		topicModel:                     repository.NewTopicModel(ctx),
		carbonCommentModel:             repository.NewCarbonCommentModel(ctx),
		carbonCommentLikeModel:         repository.NewCarbonCommentLikeRepository(ctx),
		carbonSecondHandCommodityModel: repository.NewCarbonSecondHandCommodityModel(ctx),
	}
}

func (srv *defaultCommentService) FindOne(commentId int64) (*entity.CommentIndex, error) {
	one, err := srv.commentModel.FindOne(commentId)
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (srv *defaultCommentService) FindOneQuery(data *entity.CommentIndex) (*entity.CommentIndex, error) {
	builder := srv.commentModel.RowBuilder()
	if data.Id != 0 {
		builder.Where("id = ?", data.Id)
	}
	if data.MemberId != 0 {
		builder.Where("member_id = ?", data.MemberId)
	}
	if data.RootCommentId != 0 {
		builder.Where("RootCommentId = ?", data.RootCommentId)
	}
	if data.ToCommentId != 0 {
		builder.Where("ToCommentId = ?", data.ToCommentId)
	}
	if data.ObjId != 0 {
		builder.Where("obj_id = ?", data.ObjId)
	}
	if data.Attrs != 0 {
		builder.Where("attrs = ?", data.Attrs)
	}
	builder.Where("state = ?", 0)
	query, err := srv.commentModel.FindOneQuery(builder)
	if err != nil {
		return nil, err
	}
	return query, err
}

func (srv *defaultCommentService) FindAll(data *entity.CommentIndex) ([]*entity.CommentIndex, int64, error) {
	builder := srv.commentModel.RowBuilder()
	if data.MemberId != 0 {
		builder.Where("member_id = ?", data.MemberId)
	}
	if data.RootCommentId != 0 {
		builder.Where("RootCommentId = ?", data.RootCommentId)
	}
	if data.ToCommentId != 0 {
		builder.Where("ToCommentId = ?", data.ToCommentId)
	}
	if data.ObjId != 0 {
		builder.Where("obj_id = ?", data.ObjId)
	}
	if data.Attrs != 0 {
		builder.Where("attrs = ?", data.Attrs)
	}
	builder.Where("state = ?", 0)
	all, err := srv.commentModel.FindAll(builder, "like_count DESC")
	if err != nil {
		return nil, 0, err
	}
	count, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id"))
	if err != nil {
		return nil, 0, err
	}
	return all, count, nil
}

// FindListAndChild 获取顶级评论列表及各分属下的三条子评论
func (srv *defaultCommentService) FindListAndChild(params *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error) {
	commentList := make([]*entity.CommentIndex, 0)
	var total int64
	err := app.DB.Model(&entity.CommentIndex{}).
		Preload("RootChild", func(db *gorm.DB) *gorm.DB {
			return db.Where("(select count(*) from comment_index index where index.root_comment_id = comment_index.root_comment_id and index.id <= comment_index.id) <= ?", 3).
				Order("comment_index.like_count desc").Preload("Member")
		}).
		Preload("Member").
		Where("to_comment_id = ?", 0).
		Where("obj_id = ?", params.ObjId).
		Where("state = ?", 0).
		Count(&total).
		Order("comment_index.id desc").
		Limit(limit).
		Offset(offset).
		Find(&commentList).Error
	if err != nil {
		return nil, 0, err
	}
	return commentList, total, nil
}

//FindSubList 分页获取子评论列表
func (srv *defaultCommentService) FindSubList(data *entity.CommentIndex, offSize, limit int) ([]*entity.CommentIndex, int64, error) {
	commentList := make([]*entity.CommentIndex, 0)
	var total int64
	err := srv.commentModel.RowBuilder().
		Preload("Member").
		Where("root_comment_id = ?", data.RootCommentId).
		Where("state = ?", 0).
		Count(&total).
		Order("comment_index.like_count desc, comment_index.id asc").
		Find(&commentList).Error
	if err != nil {
		return nil, 0, err
	}
	return commentList, total, nil
}

func (srv *defaultCommentService) UpdateComment(userId, commentId int64, message string) error {

	req := entity.CommentIndex{
		Id:       commentId,
		Message:  message,
		MemberId: userId,
	}
	err := srv.commentModel.Update(&req)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) DelComment(userId, commentId int64) error {
	err := srv.commentModel.Delete(commentId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) DelCommentSoft(userId, commentId int64) error {
	err := srv.commentModel.DeleteSoft(commentId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string, openId string) (entity.CommentIndex, *entity.CommentIndex, int64, error) {
	topic := srv.topicModel.FindById(topicId)

	if topic.Id == 0 {
		return entity.CommentIndex{}, &entity.CommentIndex{}, 0, errno.ErrRecordNotFound
	}

	if message != "" {
		//检查内容
		if err := validator.CheckMsgWithOpenId(openId, message); err != nil {
			app.Logger.Error(fmt.Errorf("create Comment error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "发布评论"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openId, zhuGeAttr)
			return entity.CommentIndex{}, &entity.CommentIndex{}, 0, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	comment := entity.CommentIndex{
		ObjId:         topicId,
		Message:       message,
		MemberId:      userId,
		RootCommentId: RootCommentId,
		ToCommentId:   ToCommentId,
	}
	if topic.UserId == userId {
		comment.IsAuthor = 1
	}

	_, err := srv.commentModel.Insert(&comment)
	if err != nil {
		return entity.CommentIndex{}, &entity.CommentIndex{}, 0, err
	}

	//更新count数据
	recId := topic.UserId
	toComment := &entity.CommentIndex{}
	if ToCommentId != 0 {
		//回复的评论
		toComment, err = srv.commentModel.FindOne(ToCommentId)
		if err != nil {
			return entity.CommentIndex{}, &entity.CommentIndex{}, 0, err
		}
		recId = toComment.MemberId
		err = srv.commentModel.Trans(func(tx *gorm.DB) error {
			toComment.RootCount++ //更新父评论的根评论数量
			toComment.Count++     //更新父评论的评论数量
			err = tx.Model(toComment).Select("root_count", "count").Updates(toComment).Error
			if err != nil {
				return errors.WithMessage(err, "update ToCommentIdRow:")
			}

			//本条评论
			ToCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("to_comment_id = ?", ToCommentId))
			if err != nil {
				return err
			}

			if toComment.RootCommentId != 0 {
				ToCommentId = toComment.RootCommentId
			}

			comment.RootCommentId = ToCommentId          //更新RootCommentId
			comment.Floor = int32(ToCommentIdChildCount) //更新楼层
			err = tx.Model(comment).Select("root_comment_id", "floor").Updates(comment).Error
			if err != nil {
				return errors.WithMessage(err, "update dataRow:")
			}

			//顶级评论
			if toComment.RootCommentId != 0 {
				RootCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("root_comment_id = ?", toComment.RootCommentId))
				if err != nil {
					return err
				}

				//更新顶级评论下的评论数量
				err = tx.Model(&entity.CommentIndex{}).Where("id = ?", toComment.RootCommentId).Update("count", RootCommentIdChildCount).Error
				if err != nil {
					return errors.WithMessage(err, "update RootCommentIdRow:")
				}

			}
			return nil
		})
		if err != nil {
			return entity.CommentIndex{}, &entity.CommentIndex{}, 0, err
		}

	}

	return comment, toComment, recId, nil
}

func (srv *defaultCommentService) Like(userId, commentId int64, openId string) (CommentChangeLikeResp, error) {
	comment, err := srv.commentModel.FindOne(commentId)
	if err != nil {
		return CommentChangeLikeResp{}, err
	}
	//var resp CommentChangeLikeResp

	like := srv.commentLikeModel.FindBy(repository.FindCommentLikeBy{
		CommentId: commentId,
		UserId:    userId,
	})

	var isFirst bool
	if like.Id == 0 {
		like = entity.CommentLike{
			CommentId: commentId,
			UserId:    userId,
			Status:    1,
			CreatedAt: model.Time{Time: time.Now()},
		}
		isFirst = true
	} else {
		like.UpdatedAt = model.Time{Time: time.Now()}
		like.Status = (like.Status + 1) % 2
		isFirst = false
	}

	if like.Status == 1 {
		_ = srv.AddCommentLikeCount(commentId, 1)
	} else {
		_ = srv.AddCommentLikeCount(commentId, -1)
	}

	message := comment.Message
	if len([]rune(message)) > 8 {
		message = string([]rune(message)[0:8]) + "..."
	}

	if err = srv.commentLikeModel.Save(&like); err != nil {
		return CommentChangeLikeResp{}, err
	}

	resp := CommentChangeLikeResp{
		CommentMessage: message,
		CommentId:      comment.Id,
		CommentUserId:  comment.MemberId,
		LikeStatus:     int(like.Status),
		IsFirst:        isFirst,
	}

	return resp, nil
}

func (srv *defaultCommentService) AddCommentLikeCount(commentId int64, num int) error {
	err := srv.commentModel.AddCommentLikeCount(commentId, num)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) TurnComment(params TurnCommentReq) (*APICommentResp, error) {
	kuMio, _ := util.InArray(params.Types, []int{1, 2, 3})
	mall, _ := util.InArray(params.Types, []int{10, 11, 12})
	var err error
	comment := &APICommentResp{}

	if kuMio {
		comment, err = srv.kuMioComment(params.TurnId, params.UserId)

	} else if mall {
		comment, err = srv.mallComment(params.TurnId, params.UserId)
	}

	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (srv *defaultCommentService) kuMioComment(id, userId int64) (*APICommentResp, error) {
	childList := make([]*APIComment, 0)
	comment := &APIComment{}
	commentResp := &APICommentResp{}
	commentRespChild := make([]*APICommentResp, 0)
	//root

	err := srv.ctx.DB.WithContext(srv.ctx.Context).Model(&entity.CommentIndex{}).Where("id = ?", id).First(&comment).Error
	if err != nil {
		return commentResp, err
	}

	if comment.Id == 0 {
		return commentResp, errno.ErrRecordNotFound
	}

	if comment.RootCommentId != 0 {
		id = comment.RootCommentId
	}

	//root
	err = srv.ctx.DB.WithContext(srv.ctx.Context).Model(&entity.CommentIndex{}).Where("id = ?", id).Preload("Member").First(&comment).Error
	if err != nil {
		return commentResp, err
	}

	// child
	err = srv.ctx.DB.WithContext(srv.ctx.Context).Model(&entity.CommentIndex{}).
		Preload("Member").
		Where("root_comment_id = ?", id).
		Where("state = ?", 0).
		Find(&childList).Error
	if err != nil {
		return commentResp, err
	}

	commentResp = comment.ApiComment()
	//like
	likeMap := make(map[int64]int, 0)
	commentLike := srv.commentLikeModel.GetListBy(repository.GetCommentLikeListBy{UserId: userId})
	if len(commentLike) > 0 {
		for _, item := range commentLike {
			likeMap[item.CommentId] = int(item.Status)
		}
	}

	if _, ok := likeMap[comment.Id]; ok {
		commentResp.IsLike = 1
	}

	for _, item := range childList {
		res := item.ApiComment()
		if _, ok := likeMap[comment.Id]; ok {
			res.IsLike = 1
		}
		commentRespChild = append(commentRespChild, res)
	}
	commentResp.RootChild = commentRespChild
	// obj
	obj := srv.topicModel.FindById(comment.ObjId)
	detail := Detail{
		ObjId:       obj.Id,
		ObjType:     0,
		ImageList:   obj.ImageList,
		Description: obj.Title,
	}
	commentResp.Detail = detail
	return commentResp, nil
}

func (srv *defaultCommentService) mallComment(id, userId int64) (*APICommentResp, error) {
	childList := make([]*APIComment, 0)
	comment := &APIComment{}
	commentResp := &APICommentResp{}
	commentRespChild := make([]*APICommentResp, 0)
	//root
	err := srv.ctx.DB.WithContext(srv.ctx.Context).Model(&entity.CarbonCommentIndex{}).Where("id = ?", id).First(&comment).Error
	if err != nil {
		return commentResp, err
	}

	if comment.Id == 0 {
		return commentResp, errno.ErrRecordNotFound
	}

	if comment.RootCommentId != 0 {
		id = comment.RootCommentId
	}

	//root
	err = srv.ctx.DB.WithContext(srv.ctx.Context).Model(&entity.CarbonCommentIndex{}).Where("id = ?", id).Preload("Member").First(&comment).Error
	if err != nil {
		return commentResp, err
	}

	// child
	err = srv.ctx.DB.WithContext(srv.ctx.Context).Model(&entity.CarbonCommentIndex{}).
		Preload("Member").
		Where("root_comment_id = ?", id).
		Where("state = ?", 0).
		Find(&childList).Error
	if err != nil {
		return commentResp, err
	}

	commentResp = comment.ApiComment()
	//like
	likeMap := make(map[int64]int, 0)
	commentLike := srv.carbonCommentLikeModel.GetListBy(repository.GetCommentLikeListBy{UserId: userId})
	if len(commentLike) > 0 {
		for _, item := range commentLike {
			likeMap[item.CommentId] = int(item.Status)
		}
	}

	if _, ok := likeMap[comment.Id]; ok {
		commentResp.IsLike = 1
	}

	for _, item := range childList {
		res := item.ApiComment()
		if _, ok := likeMap[comment.Id]; ok {
			res.IsLike = 1
		}
		commentRespChild = append(commentRespChild, res)
	}
	commentResp.RootChild = commentRespChild
	//obj
	obj, _ := srv.carbonSecondHandCommodityModel.FindOne(comment.ObjId)
	detail := Detail{
		ObjId:       obj.Id,
		ObjType:     1,
		ImageList:   obj.ImageList,
		Description: obj.Description,
	}
	commentResp.Detail = detail
	return commentResp, nil
}
