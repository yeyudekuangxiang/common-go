package communitymsg

import "mio/internal/pkg/model/entity"

type Topic struct {
	Event      string       `json:"event"`
	LikeStatus int          `json:"likeStatus"`
	Id         int64        `json:"id"`
	UserId     int64        `json:"userId"` // 用户
	Status     int          `json:"status"` // 状态：1 待审核 2审核失败 3已发布 4已下架
	Type       int          `json:"type"`   // 1 文章 2 活动
	Tags       []entity.Tag `json:"tags"`
}

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
