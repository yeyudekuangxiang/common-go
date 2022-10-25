package service

import (
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

func NewTopicLikeService(ctx *context.MioContext) TopicLikeService {
	return TopicLikeService{
		topicLikeModel: repository.NewTopicLikeRepository(ctx),
		topicModel:     repository.NewTopicRepository(ctx),
	}
}

type TopicLikeService struct {
	topicLikeModel repository.TopicLikeModel
	topicModel     repository.TopicModel
}

func (srv TopicLikeService) ChangeLikeStatus(topicId, userId int64, openId string) (*entity.TopicLike, int64, error) {
	topic := srv.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return nil, 0, errno.ErrCommon.WithMessage("帖子不存在")
	}
	title := topic.Title
	if len([]rune(title)) > 8 {
		title = string([]rune(topic.Title)[0:8]) + "..."
	}
	like := srv.topicLikeModel.FindBy(repository.FindTopicLikeBy{
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
		_ = srv.topicModel.AddTopicLikeCount(int64(topicId), 1)
	} else {
		_ = srv.topicModel.AddTopicLikeCount(int64(topicId), -1)
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
	return &like, point, srv.topicLikeModel.Save(&like)
}

func (srv TopicLikeService) GetLikeInfoByUser(userId int64) ([]entity.TopicLike, error) {
	list := srv.topicLikeModel.GetListBy(repository.GetTopicLikeListBy{
		UserId: userId,
		Status: 1,
	})
	if len(list) == 0 {
		return nil, errno.ErrCommon.WithMessage("未找到点赞数据")
	}
	return list, nil
}

func (srv TopicLikeService) GetOneByTopic(topicId int64, userId int64) (entity.TopicLike, error) {
	like := srv.topicLikeModel.FindBy(repository.FindTopicLikeBy{TopicId: topicId, UserId: userId})
	if like.Id == 0 {
		return entity.TopicLike{}, errno.ErrCommon.WithMessage("未找到点赞数据")
	}
	return like, nil
}
