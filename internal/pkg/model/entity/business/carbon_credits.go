package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model"
)

type CarbonCredits struct {
	ID        int64           `json:"id" gorm:"primaryKey;not null;type:serial8;comment:碳账户表"`
	BUserId   int64           `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	Credits   decimal.Decimal `json:"credits" gorm:"not null;type:decimal(20,2);comment:碳积分余额"`
	CreatedAt model.Time      `json:"createdAt" gorm:"not null;type:timestamp"`
	UpdatedAt model.Time      `json:"updatedAt" gorm:"not null;type:timestamp"`
}

func (CarbonCredits) TableName() string {
	return "business_carbon_credits"
}
