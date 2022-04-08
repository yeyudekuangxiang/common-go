package service

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	repository2 "mio/internal/pkg/repository"
	"time"
)

type TopicLikeService struct {
}

func (t TopicLikeService) ChangeLikeStatus(topicId, userId int) (*entity.TopicLike, error) {
	topic := repository2.DefaultTopicRepository.FindById(int64(topicId))
	if topic.Id == 0 {
		return nil, errors.New("帖子不存在")
	}

	r := repository2.TopicLikeRepository{DB: app.DB}
	like := r.FindBy(repository2.FindTopicLikeBy{
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
	if like.Status == 1 {
		_ = repository2.DefaultTopicRepository.AddTopicLikeCount(int64(topicId), 1)
	} else {
		_ = repository2.DefaultTopicRepository.AddTopicLikeCount(int64(topicId), -1)
	}

	return &like, r.Save(&like)
}
