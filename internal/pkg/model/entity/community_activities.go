package entity

import "mio/internal/pkg/model"

// CommunityActivities 社区用户举办的活动表
type CommunityActivities struct {
	Id             int64      `gorm:"id" json:"id"`
	Type           int        `json:"type" gorm:"type"`
	Region         string     `json:"region" gorm:"region"`
	Address        string     `json:"address" gorm:"address"`
	TagIds         string     `json:"tagIds" gorm:"tag_ids"`
	Remarks        string     `json:"remarks" gorm:"remarks"`
	Qrcode         string     `json:"qrcode" gorm:"qrcode"`
	MeetingLink    string     `json:"meetingLink" gorm:"meeting_link"`
	Contacts       string     `json:"contacts" gorm:"contacts"`
	StartTime      model.Time `json:"startTime" gorm:"start_time"`
	EndTime        model.Time `json:"endTime" gorm:"end_time"`
	SignupDeadline model.Time `json:"signupDeadline" gorm:"signup_deadline"`
	//数据库没有的字段
	Status       int `json:"status,omitempty" gorm:"-"`
	SignupStatus int `json:"signupStatus" gorm:"-"`
}

func (CommunityActivities) TableName() string {
	return "community_activities"
}
