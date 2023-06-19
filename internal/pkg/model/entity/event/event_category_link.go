package event

type EventCategoryLink struct {
	ID              int64  `json:"id" gorm:"primaryKey;type:serial4;not null;comment:公益活动分类链接表"` //分类id
	EventCategoryId int64  `json:"eventCategoryId"`
	Link            string `json:"link" gorm:"type:varchar(255);not null;comment:链接"` //分类标识
	Display         int    `json:"display"`
}

func (EventCategoryLink) TableName() string {
	return "event_category_link"
}
