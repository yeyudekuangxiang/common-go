package service

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"strconv"
	"time"
)

var DefaultTopicLikeService = NewDefaultTopicLikeService()

func NewDefaultTopicLikeService() TopicLikeService {
	return TopicLikeService{
		repo: repository.DefaultTopicLikeRepository,
	}
}

type TopicLikeService struct {
	repo repository.TopicLikeRepository
}

func (t TopicLikeService) ChangeLikeStatus(topicId, userId int, openId string) (*entity.TopicLike, error) {
	topic := repository.DefaultTopicRepository.FindById(int64(topicId))
	if topic.Id == 0 {
		return nil, errors.New("帖子不存在")
	}
	title := topic.Title
	if len(topic.Title) > 8 {
		title = topic.Title[0:8] + "..."
	}
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
	if like.Status == 1 {
		_ = repository.DefaultTopicRepository.AddTopicLikeCount(int64(topicId), 1)
		pointService := NewPointService(context.NewMioContext())
		_, _ = pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       openId,
			Type:         entity.POINT_LIKE,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_LIKE]),
			AdminId:      0,
			Note:         "为文章 \"" + title + "\" 点赞",
			AdditionInfo: strconv.FormatInt(topic.Id, 10),
		})
	} else {
		_ = repository.DefaultTopicRepository.AddTopicLikeCount(int64(topicId), -1)
	}

	return &like, r.Save(&like)
}

func (t TopicLikeService) GetLikeInfoByUser(userId int64) ([]entity.TopicLike, error) {
	list := t.repo.GetListBy(repository.GetTopicLikeListBy{UserId: userId})
	if len(list) == 0 {
		return nil, errors.New("未找到数据")
	}
	return list, nil
}
