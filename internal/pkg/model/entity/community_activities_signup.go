package entity

import (
	"time"
)

// CommunityActivitiesSignup 社区活动 报名表
type CommunityActivitiesSignup struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TopicId      int64     `json:"topicId"`
	UserId       int64     `json:"userId"`
	RealName     string    `json:"realName"`
	Phone        string    `json:"phone"`
	Gender       int       `json:"gender"`
	Age          int       `json:"age"`
	Wechat       string    `json:"wechat"`
	City         string    `json:"city"`
	Remarks      string    `json:"remarks"`
	SignupTime   time.Time `json:"signupTime"`
	CancelTime   time.Time `json:"cancelTime,omitempty"`
	SignupStatus int       `json:"signupStatus"`
}

func (CommunityActivitiesSignup) TableName() string {
	return "community_activities_signup"
}

type APIActivitiesSignup struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"userId"`
	RealName     string    `json:"realName"`
	Phone        string    `json:"phone"`
	Gender       int       `json:"gender"`
	Age          int       `json:"age"`
	Wechat       string    `json:"wechat"`
	City         string    `json:"city"`
	Remarks      string    `json:"remarks"`
	SignupTime   time.Time `json:"signupTime"`
	CancelTime   time.Time `json:"cancelTime"`
	SignupStatus int       `json:"signupStatus"`
	TopicId      int64     `json:"topicId"`
	Topic        APITopic  `json:"topic" gorm:"foreignKey:TopicId"`
}

type APISignupList struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"userId"`
	User         ShortUser `json:"user"`
	SignupTime   time.Time `json:"signupTime"`
	SignupStatus int       `json:"signupStatus"`
}
