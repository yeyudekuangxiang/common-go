package business

import "mio/internal/pkg/model"

type CompanyCarbonScene struct {
	ID               int64            `json:"id" gorm:"primaryKey;not null;type:serial8;comment:企业低碳场景表"`
	CarbonSceneId    int              `json:"carbonSceneId" gorm:"not null;type:int4;comment:低碳场景表主键"`
	BCompanyId       int              `json:"bCompanyId" gorm:"not null;type:int4;comment:所属企业版企业表主键id"`
	Sort             int              `json:"sort" gorm:"not null;type:int4;default:999;comment:排序 从小到大排序"`
	Status           int8             `json:"status" gorm:"not null;type:int2;comment:上线状态 1上线 2下线"`
	PointRateSetting PointRateSetting `json:"pointSetting" gorm:"not null;type:varchar(500);comment:每个单位或者每次获取的积分值"`
	MaxCount         int              `json:"maxCount" gorm:"not null;type:int4;comment:单日最大积分数量"`
	MaxPoint         int              `json:"maxPoint" gorm:"not null;type:int4;comment:单日最大积分值"`
	CreatedAt        model.Time       `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt        model.Time       `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (CompanyCarbonScene) TableName() string {
	return "business_company_carbon_scene"
}
