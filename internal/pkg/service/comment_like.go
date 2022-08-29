package service

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"time"
)

var DefaultCommentLikeService = CommentLikeService{repo: repository.DefaultCommentLikeRepository}

type CommentLikeService struct {
	repo repository.CommentLikeRepository
}

func (t CommentLikeService) Like(userId, commentId int64) (*entity.CommentLike, error) {
	like := t.repo.FindBy(repository.FindCommentLikeBy{
		CommentId: commentId,
		UserId:    userId,
	})
	if like.Id == 0 {
		like.Status = 1
		like.CommentId = commentId
		like.UserId = userId
		like.CreatedAt = model.Time{Time: time.Now()}
	} else {
		like.UpdatedAt = model.Time{Time: time.Now()}
		like.Status = (like.Status + 1) % 2
	}
	if like.Status == 1 {
		_ = DefaultCommentService.AddTopicLikeCount(commentId, 1)
	} else {
		_ = DefaultCommentService.AddTopicLikeCount(commentId, -1)
	}

	return &like, t.repo.Save(&like)
}
