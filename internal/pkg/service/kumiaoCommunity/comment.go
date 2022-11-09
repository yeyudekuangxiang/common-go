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
		CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string, openId string) (entity.CommentIndex, int64, error)
		UpdateComment(userId, commentId int64, message string) error
		DelComment(userId, commentId int64) error
		DelCommentSoft(userId, commentId int64) error
		Like(userId, commentId int64, openId string) (CommentChangeLikeResp, error)
		AddCommentLikeCount(commentId int64, num int) error
	}
)

type defaultCommentService struct {
	ctx              *context.MioContext
	commentModel     repository.CommentModel
	commentLikeModel repository.CommentLikeModel
	topicModel       repository.TopicModel
}

func NewCommentService(ctx *context.MioContext) CommentService {
	return &defaultCommentService{
		ctx:              ctx,
		commentModel:     repository.NewCommentModel(ctx),
		commentLikeModel: repository.NewCommentLikeRepository(ctx),
		topicModel:       repository.NewTopicModel(ctx),
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
		Order("comment_index.like_count desc,comment_index.id asc").
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

func (srv *defaultCommentService) CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string, openId string) (entity.CommentIndex, int64, error) {
	topic := srv.topicModel.FindById(topicId)

	if topic.Id == 0 {
		return entity.CommentIndex{}, 0, errno.ErrRecordNotFound
	}

	if message != "" {
		//检查内容
		if err := validator.CheckMsgWithOpenId(openId, message); err != nil {
			app.Logger.Error(fmt.Errorf("create Comment error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "发布评论"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openId, zhuGeAttr)
			return entity.CommentIndex{}, 0, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	comment := entity.CommentIndex{
		ObjId:         topicId,
		Message:       message,
		MemberId:      userId,
		RootCommentId: RootCommentId,
		ToCommentId:   ToCommentId,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
	if topic.UserId == userId {
		comment.IsAuthor = 1
	}

	_, err := srv.commentModel.Insert(&comment)
	if err != nil {
		return entity.CommentIndex{}, 0, err
	}

	//更新topic
	//app.DB.Model(&topic).Update("updated_at", model.Time{Time: time.Now()})
	//更新count数据
	recId := topic.UserId
	if ToCommentId != 0 {
		//回复的评论
		ToComment, err := srv.commentModel.FindOne(ToCommentId)
		if err != nil {
			return entity.CommentIndex{}, 0, err
		}
		recId = ToComment.MemberId
		err = srv.commentModel.Trans(func(tx *gorm.DB) error {
			ToComment.RootCount++ //更新父评论的根评论数量
			ToComment.Count++     //更新父评论的评论数量
			err = tx.Model(ToComment).Select("root_count", "count").Updates(ToComment).Error
			if err != nil {
				return errors.WithMessage(err, "update ToCommentIdRow:")
			}

			//本条评论
			ToCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("to_comment_id = ?", ToCommentId))
			if err != nil {
				return err
			}

			if ToComment.RootCommentId != 0 {
				ToCommentId = ToComment.RootCommentId
			}

			comment.RootCommentId = ToCommentId          //更新RootCommentId
			comment.Floor = int32(ToCommentIdChildCount) //更新楼层
			err = tx.Model(comment).Select("root_comment_id", "floor").Updates(comment).Error
			if err != nil {
				return errors.WithMessage(err, "update dataRow:")
			}

			//顶级评论
			if ToComment.RootCommentId != 0 {
				RootCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("root_comment_id = ?", ToComment.RootCommentId))
				if err != nil {
					return err
				}

				//更新顶级评论下的评论数量
				err = tx.Model(&entity.CommentIndex{}).Where("id = ?", ToComment.RootCommentId).Update("count", RootCommentIdChildCount).Error
				if err != nil {
					return errors.WithMessage(err, "update RootCommentIdRow:")
				}

			}
			return nil
		})
		if err != nil {
			return entity.CommentIndex{}, 0, err
		}

	}

	return comment, recId, nil
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
