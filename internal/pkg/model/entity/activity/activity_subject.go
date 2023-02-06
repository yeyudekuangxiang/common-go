package activity

import "time"

type Subject struct {
	Id            int64     `json:"id"`
	ActivityId    int64     `json:"activityId"`
	Title         string    `json:"title"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	Status        int       `json:"status"`
	UperRiskLimit int       `json:"uperRiskLimit"`
	CreatedAt     time.Time `json:"createdAt"`
	Creator       string    `json:"creator"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Updater       string    `json:"updater"`
}

func (Subject) TableName() string {
	return "activity_subject"
}
