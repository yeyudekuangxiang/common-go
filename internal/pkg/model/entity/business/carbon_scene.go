package business

import (
	"mio/internal/pkg/model"
)

type CarbonScene struct {
	ID           int        `json:"id" gorm:"primaryKey;not null;type:serial4;comment:低碳场景表"`
	Type         CarbonType `json:"type" gorm:"not null;type:varchar(50);comment:低碳场景类型"`
	PointSetting int        `json:"pointSetting" gorm:"not null;type:int4;comment:每个单位或者每次获取的积分值"`
	Title        string     `json:"title" gorm:"not null;type:varchar(100);comment:场景标题"`
	Icon         string     `json:"icon" gorm:"not null;type:varchar(500);comment:场景icon链接"`
	Desc         string     `json:"desc" gorm:"not null;type:varchar(255);comment:场景描述"`
	MaxCount     int        `json:"maxCount" gorm:"not null;type:int4;comment:每日最多获取积分次数"`
	MaxPoint     int        `json:"maxPoint" gorm:"not null;type:int4;comment:每日最多获取积分值"`
	CreatedAt    model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt    model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (CarbonScene) TableName() string {
	return "business_carbon_scene"
}
