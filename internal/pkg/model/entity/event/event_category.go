package event

import "mio/internal/pkg/model/entity"

var OrderByEventCategorySortDesc entity.OrderBy = "order_by_event_category_sort_desc"

type EventCategory struct {
	ID              int    `json:"id" gorm:"primaryKey;type:serial4;not null;comment:公益活动分类表"`         //分类id
	EventCategoryId string `json:"eventCategoryId" gorm:"type:varchar(255);not null;comment:分类标识"`     //分类标识
	Title           string `json:"title" gorm:"type:varchar(255);not null;comment:分类名称"`               //分类标识
	Active          bool   `json:"active" gorm:"type:bool;not null;default:false;comment:是否上线"`        //是否上线
	ImageUrl        string `json:"imageUrl" gorm:"type:varchar(500);not null;default:'';comment:分类主图"` //分类主图 尺寸 1372 × 480
}

func (EventCategory) TableName() string {
	return "event_category"
}
