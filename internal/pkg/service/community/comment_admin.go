package community

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
)

type (
	CommentAdminService interface {
		DelCommentSoft(commentId int64, reason string) (entity.CommentIndex, error)
		CommentList(content string, userId int64, topicId int64, limit, offset int) ([]*entity.CommentIndex, int64, error)
	}
)

type defaultCommentAdminService struct {
	ctx          *context.MioContext
	commentModel repository.CommentModel
}

func (d defaultCommentAdminService) CommentList(content string, userId int64, objId int64, limit, offset int) ([]*entity.CommentIndex, int64, error) {
	builder := d.commentModel.RowBuilder().Preload("Member")
	if content != "" {
		builder.Where("message like ?", "%"+content+"%")
	}
	if userId != 0 {
		builder.Where("member_id = ?", userId)
		builder.Preload("Member")
	}
	if objId != 0 {
		builder.Where("obj_id = ?", objId)
	}
	count, err := d.commentModel.FindCount(builder)
	if err != nil {
		return nil, 0, err
	}
	builder.Limit(limit).Offset(offset)
	all, err := d.commentModel.FindAll(builder, "id desc, like_count desc")
	if err != nil {
		if err == entity.ErrNotFount {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return all, count, nil
}

func (d defaultCommentAdminService) DelCommentSoft(commentId int64, reason string) (entity.CommentIndex, error) {
	builder := d.commentModel.RowBuilder()
	builder.Where("id = ?", commentId).Where("state = ?", 0)
	comment, err := d.commentModel.FindOneQuery(builder)
	if err != nil {
		if err == entity.ErrNotFount {
			return entity.CommentIndex{}, errno.ErrCommon.WithMessage("该评论不存在")
		}
		return entity.CommentIndex{}, err
	}
	err = d.commentModel.RowBuilder().Model(comment).Updates(entity.CommentIndex{State: 1, DelReason: reason}).Error
	if err != nil {
		return entity.CommentIndex{}, err
	}
	return *comment, nil
}

func NewCommentAdminService(ctx *context.MioContext) CommentAdminService {
	return &defaultCommentAdminService{
		ctx:          ctx,
		commentModel: repository.NewCommentModel(ctx),
	}
}
