package business

import "mio/internal/pkg/model"

type PointLimitLog struct {
	ID           int64      `json:"id" gorm:"primaryKey;not null;type:serial8;comment:积分获取次数限制表"`
	Type         string     `json:"type" gorm:"not null;type:varchar(20);comment:积分类型"`
	BUserId      int64      `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	CurrentCount int        `json:"currentCount" gorm:"not null;type:int4;comment:已获取积分次数"`
	CreatedAt    model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt    model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (PointLimitLog) TableName() string {
	return "business_point_limit_log"
}
