package kumiaoCommunity

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/errno"
	"time"
)

func NewTopicLikeService(ctx *context.MioContext) TopicLikeService {
	return TopicLikeService{
		topicLikeModel: repository.NewTopicLikeRepository(ctx),
		topicModel:     repository.NewTopicModel(ctx),
	}
}

type TopicLikeService struct {
	topicLikeModel repository.TopicLikeModel
	topicModel     repository.TopicModel
}

func (srv TopicLikeService) ChangeLikeStatus(topicId, userId int64, openId string) (TopicChangeLikeResp, error) {
	topic := srv.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return TopicChangeLikeResp{}, errno.ErrCommon.WithMessage("帖子不存在")
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
		_ = srv.topicModel.AddTopicLikeCount(topicId, 1)
	} else {
		_ = srv.topicModel.AddTopicLikeCount(topicId, -1)
	}

	err := srv.topicLikeModel.Save(&like)
	if err != nil {
		return TopicChangeLikeResp{}, err
	}

	return TopicChangeLikeResp{
		TopicTitle:  topic.Title,
		TopicId:     topic.Id,
		TopicUserId: topic.User.ID,
		LikeStatus:  int(like.Status),
		IsFirst:     isFirst,
	}, nil
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