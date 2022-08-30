package service

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"time"
)

var DefaultCommentService = NewCommentService(repository.DefaultCommentRepository)

type (
	CommentService interface {
		FindOne(commentId int64) (*entity.CommentIndex, error)
		FindOneQuery(data *entity.CommentIndex) (*entity.CommentIndex, error)
		FindAll(data *entity.CommentIndex) ([]*entity.CommentIndex, int64, error)
		FindSubList(data *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error)
		FindListAndChild(data *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error)
		CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string) (entity.CommentIndex, error)
		UpdateComment(userId, commentId int64, message string) error
		DelComment(userId, commentId int64) error
		DelCommentSoft(userId, commentId int64) error
		Like(userId, commentId int64) error
		AddTopicLikeCount(commentId int64, num int) error
	}
)

type defaultCommentService struct {
	commentModel repository.CommentModel
}

func NewCommentService(model repository.CommentModel) CommentService {
	return &defaultCommentService{
		commentModel: model,
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
	if data.ID != 0 {
		builder.Where("id = ?", data.ID)
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
		Order("like_count desc, id asc").
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
		Order("like_count desc, id asc").
		Find(&commentList).Error
	if err != nil {
		return nil, 0, err
	}
	return commentList, total, nil
}

func (srv *defaultCommentService) UpdateComment(userId, commentId int64, message string) error {
	req := entity.CommentIndex{
		ID:       commentId,
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

func (srv *defaultCommentService) CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string) (entity.CommentIndex, error) {
	comment := entity.CommentIndex{
		ObjId:         topicId,
		Message:       message,
		MemberId:      userId,
		RootCommentId: RootCommentId,
		ToCommentId:   ToCommentId,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
	//to user info
	//var toUser *entity.User
	//if ToCommentId != 0 {
	//	toComment, _ := srv.FindOne(ToCommentId)
	//	toUser, _ = DefaultUserService.GetUserById(toComment.MemberId)
	//	comment.ToNickName = toUser.Nickname
	//}
	_, err := srv.commentModel.Insert(&comment)
	if err != nil {
		return comment, err
	}
	//更新count数据
	if ToCommentId != 0 {
		err = srv.commentModel.Trans(func(tx *gorm.DB) error {
			//回复的评论
			ToCommentIdRow, err := srv.commentModel.FindOne(ToCommentId)
			if err != nil {
				return err
			}
			ToCommentIdRow.RootCount++ //更新父评论的根评论数量
			ToCommentIdRow.Count++     //更新父评论的评论数量
			err = app.DB.Model(ToCommentIdRow).Select("RootCommentId_count", "count").Updates(ToCommentIdRow).Error
			if err != nil {
				return errors.WithMessage(err, "update ToCommentIdRow:")
			}
			//本条评论
			ToCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("ToCommentId = ?", ToCommentId))
			if err != nil {
				return err
			}
			if ToCommentIdRow.RootCommentId != 0 {
				ToCommentId = ToCommentIdRow.RootCommentId
			}
			comment.RootCommentId = ToCommentId          //更新RootCommentId
			comment.Floor = int32(ToCommentIdChildCount) //更新楼层
			err = app.DB.Model(comment).Select("RootCommentId", "floor").Updates(comment).Error
			if err != nil {
				return errors.WithMessage(err, "update dataRow:")
			}

			//顶级评论
			if ToCommentIdRow.RootCommentId != 0 {
				RootCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("RootCommentId = ?", ToCommentIdRow.RootCommentId))
				if err != nil {
					return err
				}
				//更新顶级评论下的评论数量
				err = app.DB.Model(&entity.CommentIndex{}).Where("id = ?", ToCommentIdRow.RootCommentId).Update("count", RootCommentIdChildCount).Error
				if err != nil {
					return errors.WithMessage(err, "update RootCommentIdRow:")
				}
			}
			return nil
		})
		if err != nil {
			return comment, err
		}
	}
	return comment, nil
}

func (srv *defaultCommentService) Like(userId, commentId int64) error {
	_, err := srv.commentModel.FindOne(commentId)
	if err != nil {
		if err == entity.ErrNotFount {
			return nil
		}
		return err
	}
	_, err = DefaultCommentLikeService.Like(userId, commentId)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) AddTopicLikeCount(commentId int64, num int) error {
	err := srv.commentModel.AddTopicLikeCount(commentId, num)
	if err != nil {
		return err
	}
	return nil
}
