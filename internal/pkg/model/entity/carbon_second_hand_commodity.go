package entity

import (
	"mio/internal/pkg/model"
	"time"
)

type CarbonSecondHandCommodity struct {
	Id           int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64   `gorm:"type:int8;not null" json:"userId""`
	OpenId       string  `json:"openId"`
	Avatar       string  `json:"avatar"`
	Nickname     string  `json:"nickname"`
	ImageList    string  `json:"imageList"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Address      string  `json:"address"`
	Description  string  `json:"description"`
	CommentCount int64   `json:"commentCount"`
	LikeCount    int64   `json:"likeCount"`
	ViewCount    int64   `json:"viewCount"`

	Pubdate       time.Time `json:"pubdate"`
	AuditReason   string    `json:"auditReason"`
	IsRecommend   int       `json:"isRecommend"`
	RecommendTime time.Time `json:"recommendTime"`
	LastEditTime  time.Time `json:"lastEditTime"`
	Partners      int       `json:"partners"`

	State int8 `gorm:"type:int2;not null" json:"status" json:"state,omitempty"`

	CreatedAt model.Time `json:"createdAt" json:"createdAt"`
	UpdatedAt model.Time `json:"updatedAt" json:"updatedAt"`
}

func (CarbonSecondHandCommodity) TableName() string {
	return "carbon_second_hand_commodity"
}
