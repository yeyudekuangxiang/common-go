package event

type EventDetail struct {
	ID      int    `json:"id" gorm:"primaryKey;type:serial4;not null;comment:公益活动详情表"`
	EventId string `json:"eventId" gorm:"type:varchar(255);not null;comment:公益活动标识"`
	Content string `json:"content" gorm:"type:text;not null;comment:内容"`
}

func (EventDetail) TableName() string {
	return "event_detail"
}
