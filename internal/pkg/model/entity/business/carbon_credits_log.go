package business

import "mio/internal/pkg/model"

type CarbonCreditsLog struct {
	ID        int64      `json:"id" gorm:"primaryKey;not null;type:serial8;comment:碳积分变动表"`
	BUserId   int64      `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	Type      string     `json:"type" gorm:"not null;type:varchar(100);comment:减碳场景类型"`
	Value     float64    `json:"value" gorm:"not null;type:decimal(10,2);comment:获取到的碳积分数量"`
	Info      string     `json:"info" gorm:"not null;type:varchar(1000);default:'';comment:附带信息 json object 同一个type的info格式必须统一"`
	CreatedAt model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (CarbonCreditsLog) TableName() string {
	return "business_carbon_credits_log"
}
