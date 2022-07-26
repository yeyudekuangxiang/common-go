package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model"
)

type CarbonCreditsLimitLog struct {
	ID           int64           `json:"id" gorm:"primaryKey;not null;type:serial8;comment:碳积分获取限制表"`
	Type         CarbonType      `json:"type" gorm:"not null;type:varchar(50);comment:碳积分类型"`
	BUserId      int64           `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	CurrentCount int             `json:"currentCount" gorm:"not null;type:int4;comment:当天已获取碳积分次数"`
	CurrentValue decimal.Decimal `json:"currentValue" gorm:"not null;type:decimal(10,2);comment:当天已获得碳积分值"`
	TimePoint    model.Time      `json:"timePoint" gorm:"not null;type:timestamptz;comment:时间点 2006-01-02 00:00:00"`
	CreatedAt    model.Time      `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt    model.Time      `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (CarbonCreditsLimitLog) TableName() string {
	return "business_carbon_credits_limit_log"
}
