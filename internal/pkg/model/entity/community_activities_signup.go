package entity

import (
	"mio/internal/pkg/model"
	"time"
)

// CommunityActivitiesSignup 社区活动 报名表
type CommunityActivitiesSignup struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TopicId      int64     `json:"topicId"`
	UserId       int64     `json:"userId"`
	SignupInfo   string    `json:"signupInfo"`
	SignupTime   time.Time `json:"signupTime"`
	CancelTime   time.Time `json:"cancelTime,omitempty"`
	SignupStatus int       `json:"signupStatus"`
}

// CommunityActivitiesSignup 社区活动 报名表

type CommunityActivitiesSignupV2 struct {
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
	Id           int64              `json:"id"`
	SignupInfo   string             `json:"signupInfo"`
	SignupTime   model.Time         `json:"signupTime"`
	CancelTime   model.Time         `json:"cancelTime"`
	SignupStatus int                `json:"signupStatus"`
	UserId       int64              `json:"userId"`
	User         ShortUser          `json:"user" gorm:"foreignKey:UserId"`
	TopicId      int64              `json:"topicId"`
	Topic        APITopicActivities `json:"topic,omitempty" gorm:"foreignKey:TopicId"`
}

type APISignupList struct {
	Id         int64  `json:"id"`
	SignupInfo string `json:"signupInfo"`
	//关联数据
	UserId       int64      `json:"userId"`
	User         ShortUser  `json:"user"`
	SignupTime   model.Time `json:"signupTime"`
	SignupStatus int        `json:"signupStatus"`
}

type APIListCount struct {
	TopicId     int64 `json:"topicId"`
	NumOfSignup int64 `json:"numOfSignup"`
}
