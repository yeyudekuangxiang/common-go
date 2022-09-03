package service

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"strconv"
	"time"
)

var DefaultCommentLikeService = CommentLikeService{repo: repository.DefaultCommentLikeRepository}

type CommentLikeService struct {
	repo repository.CommentLikeRepository
}

func (t CommentLikeService) Like(userId, commentId int64, message, openId string) (*entity.CommentLike, int64, error) {
	like := t.repo.FindBy(repository.FindCommentLikeBy{
		CommentId: commentId,
		UserId:    userId,
	})
	var isFirst bool
	if like.Id == 0 {
		like = entity.CommentLike{
			CommentId: commentId,
			UserId:    userId,
			Status:    1,
			CreatedAt: model.Time{Time: time.Now()},
		}
		isFirst = true
	} else {
		like.UpdatedAt = model.Time{Time: time.Now()}
		like.Status = (like.Status + 1) % 2
		isFirst = false
	}
	if like.Status == 1 {
		_ = DefaultCommentService.AddTopicLikeCount(commentId, 1)
	} else {
		_ = DefaultCommentService.AddTopicLikeCount(commentId, -1)
	}
	if len(message) > 8 {
		message = string([]rune(message)[0:8]) + "..."
	}
	//发放积分
	var point int64
	if like.Status == 1 && isFirst {
		pointService := NewPointService(context.NewMioContext())
		_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       openId,
			Type:         entity.POINT_LIKE,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_LIKE]),
			AdminId:      0,
			Note:         "评论 \"" + message + "\" 点赞",
			AdditionInfo: strconv.FormatInt(commentId, 10) + "#" + strconv.FormatInt(like.Id, 10),
		})
		if err == nil {
			point = int64(entity.PointCollectValueMap[entity.POINT_LIKE])
		}
	}
	if err := t.repo.Save(&like); err != nil {
		return nil, 0, err
	}
	return &like, point, nil
}

func (t CommentLikeService) GetLikeInfoByUser(userId int64) []entity.CommentLike {
	return t.repo.GetListBy(repository.GetCommentLikeListBy{UserId: userId})
}
