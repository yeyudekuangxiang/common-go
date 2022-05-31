package business

import "mio/internal/pkg/model"

type CompanyCarbonScene struct {
	ID             int64      `json:"id" gorm:"primaryKey;not null;type:serial8;comment:企业低碳场景表"`
	BCarbonSceneId int        `json:"bCarbonSceneId" gorm:"not null;type:int4;comment:低碳场景表主键"`
	Sort           int        `json:"sort" gorm:"not null;type:int4;default:999;comment:排序 从小到大排序"`
	Status         int8       `json:"status" gorm:"not null;type:int2;comment:上线状态 1上线 2下线"`
	PointSetting   string     `json:"pointSetting" gorm:"not null;type:varchar(500);comment:企业自定义积分规则"`
	MaxCount       int        `json:"maxCount" gorm:"not null;type:int4;comment:企业自定义积分数量限制"`
	CreatedAt      model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt      model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (CompanyCarbonScene) TableName() string {
	return "business_company_carbon_scene"
}
