package community

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultCommentLikeService = CommentLikeService{}

type (
	CommentLikeModel interface {
		GetLikeInfoByUser(userId int64) []entity.CommentLike
	}

	defaultCommentLikeService struct {
		commentLikeModel repository.CommentLikeModel
	}
)

func NewCommentLikeService(ctx *context.MioContext) CommentLikeModel {
	return &defaultCommentLikeService{
		commentLikeModel: repository.NewCommentLikeRepository(ctx),
	}
}

type CommentLikeService struct {
	repo repository.CommentLikeRepository
}

func (d defaultCommentLikeService) GetLikeInfoByUser(userId int64) []entity.CommentLike {
	return d.commentLikeModel.GetListBy(repository.GetCommentLikeListBy{UserId: userId})
}
