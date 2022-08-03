package service

import (
	"errors"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultCommentAdminService = NewCommentAdminService(repository.DefaultCommentRepository)

type (
	CommentAdminService interface {
		DelCommentSoft(commentId int64) error
		CommentList(content string, userId, topicId int64, limit, offset int) ([]*entity.CommentIndex, int64, error)
	}
)

type defaultCommentAdminService struct {
	commentModel repository.CommentModel
}

func (d defaultCommentAdminService) CommentList(content string, userId, objId int64, limit, offset int) ([]*entity.CommentIndex, int64, error) {
	builder := d.commentModel.RowBuilder()
	builder.Limit(limit).Offset(offset).Preload("Member")
	if content != "" {
		builder.Where("message like ?", "%"+content+"%")
	}
	if userId != 0 {
		builder.Preload("Member", "id = ?", userId)
	}
	if objId != 0 {
		builder.Where("obj_id = ?", objId)
	}
	all, err := d.commentModel.FindAll(builder, "id desc, like_count desc")
	if err != nil {
		if err == entity.ErrNotFount {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	count, err := d.commentModel.FindCount(d.commentModel.CountBuilder("id"))
	if err != nil {
		return nil, 0, err
	}
	return all, count, nil
}

func (d defaultCommentAdminService) DelCommentSoft(commentId int64) error {
	one, err := d.commentModel.FindOne(commentId)
	if err != nil {
		if err == entity.ErrNotFount {
			return errors.New("该评论不存在")
		}
		return err
	}
	err = d.commentModel.RowBuilder().Model(one).Update("status", 4).Error
	if err != nil {
		return err
	}
	return nil
}

func NewCommentAdminService(model repository.CommentModel) CommentAdminService {
	return &defaultCommentAdminService{
		commentModel: model,
	}
}
