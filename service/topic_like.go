package service

import (
	"mio/core/app"
	"mio/model"
	"mio/model/entity"
	"mio/repository"
	"time"
)

type TopicLikeService struct {
}

func (t TopicLikeService) ChangeLikeStatus(topicId, userId int) (*entity.TopicLike, error) {
	r := repository.TopicLikeRepository{DB: app.DB}
	like := r.FindBy(repository.FindTopicLikeBy{
		TopicId: topicId,
		UserId:  userId,
	})
	if like.Id == 0 {
		like.Status = 1
		like.TopicId = topicId
		like.UserId = userId
		like.CreatedAt = model.Time{Time: time.Now()}
	} else {
		like.UpdatedAt = model.Time{Time: time.Now()}
		like.Status = (like.Status + 1) % 2
	}
	return &like, r.Save(&like)
}
