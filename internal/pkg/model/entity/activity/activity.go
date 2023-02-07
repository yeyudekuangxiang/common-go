package activity

import "time"

type Activity struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Code      int       `json:"code"`
	//StartTime time.Time `json:"startTime"`
	//EndTime   time.Time `json:"endTime"`
	Channel int       `json:"channel"`
	Subject []Subject `json:"subject" gorm:"foreignKey:ActivityId"`
}

func (Activity) TableName() string {
	return "activity"
}
