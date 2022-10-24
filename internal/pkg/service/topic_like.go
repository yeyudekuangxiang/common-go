package service

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultTopicLikeService = NewDefaultTopicLikeService()

func NewDefaultTopicLikeService() TopicLikeService {
	return TopicLikeService{
		repo:       repository.DefaultTopicLikeRepository,
		topicModel: repository.NewTopicRepository(context.NewMioContext()),
	}
}

type TopicLikeService struct {
	repo       repository.TopicLikeRepository
	topicModel repository.TopicModel
}

func (t TopicLikeService) ChangeLikeStatus(topicId, userId int, openId string) (*entity.TopicLike, int64, error) {
	topic := t.topicModel.FindById(int64(topicId))
	if topic.Id == 0 {
		return nil, 0, errno.ErrCommon.WithMessage("帖子不存在")
	}
	title := topic.Title
	if len([]rune(title)) > 8 {
		title = string([]rune(topic.Title)[0:8]) + "..."
	}
	r := repository.TopicLikeRepository{DB: app.DB}
	like := r.FindBy(repository.FindTopicLikeBy{
		TopicId: topicId,
		UserId:  userId,
	})
	var isFirst bool
	if like.Id == 0 {
		like.Status = 1
		like.TopicId = topicId
		like.UserId = userId
		like.CreatedAt = model.Time{Time: time.Now()}
		isFirst = true
	} else {
		like.UpdatedAt = model.Time{Time: time.Now()}
		like.Status = (like.Status + 1) % 2
		isFirst = false
	}
	if like.Status == 1 {
		_ = t.topicModel.AddTopicLikeCount(int64(topicId), 1)
	} else {
		_ = t.topicModel.AddTopicLikeCount(int64(topicId), -1)
	}
	var point int64
	if like.Status == 1 && isFirst == true {
		pointService := NewPointService(context.NewMioContext())
		_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       openId,
			Type:         entity.POINT_LIKE,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_LIKE]),
			AdminId:      0,
			Note:         "为文章 \"" + title + "\" 点赞",
			AdditionInfo: strconv.FormatInt(topic.Id, 10),
		})
		if err == nil {
			point = int64(entity.PointCollectValueMap[entity.POINT_LIKE])
		}
	}
	return &like, point, r.Save(&like)
}

func (t TopicLikeService) GetLikeInfoByUser(userId int64) ([]entity.TopicLike, error) {
	list := t.repo.GetListBy(repository.GetTopicLikeListBy{UserId: userId})
	if len(list) == 0 {
		return nil, errno.ErrCommon.WithMessage("未找到点赞数据")
	}
	return list, nil
}

func (t TopicLikeService) GetOneByTopic(topicId int64, userId int64) (entity.TopicLike, error) {
	like := t.repo.FindBy(repository.FindTopicLikeBy{TopicId: int(topicId), UserId: int(userId)})
	if like.Id == 0 {
		return entity.TopicLike{}, errno.ErrCommon.WithMessage("未找到点赞数据")
	}
	return like, nil
}
