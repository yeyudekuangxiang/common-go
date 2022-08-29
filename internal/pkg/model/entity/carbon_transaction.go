package entity

import (
	"time"
)

const (
	OrderByCarbonTranDayVDate OrderBy = "order_by_carbon_day_v_date_desc"
)

type CarbonTransactionType string

// Text 获取积分类型的中文名称(给用户看的)
func (p CarbonTransactionType) Text() string {
	switch p {
	case CARBON_STEP:
		return "步行"
	case CARBON_COFFEE_CUP:
		return "自带咖啡杯"
	case CARBON_BIKE_RIDE:
		return "骑行"
	case CARBON_ECAR:
		return "电动车"
	}
	return "未知积分"
}

func (p CarbonTransactionType) Cover() string {
	switch p {
	case CARBON_STEP:
		return "https://resources.miotech.com/static/mp2c/images/mp2c2.0/assets/classify_foot.png"
	case CARBON_COFFEE_CUP:
		return "https://resources.miotech.com/static/mp2c/images/mp2c2.0/assets/classify_cup.png"
	case CARBON_BIKE_RIDE:
		return "https://resources.miotech.com/static/mp2c/images/mp2c2.0/assets/classify_riding.png"
	case CARBON_ECAR:
		return "https://resources.miotech.com/static/mp2c/images/mp2c2.0/assets/classify_newenergy.png"
	}
	return ""
}

// RealText 获取积分类型的中文名称(给管理员看的)
func (p CarbonTransactionType) RealText() string {
	switch p {
	case CARBON_STEP:
		return "步行"
	case CARBON_COFFEE_CUP:
		return "自带咖啡杯"
	case CARBON_BIKE_RIDE:
		return "骑行"
	case CARBON_ECAR:
		return "电动车"
	}
	return p.Text()
}

const (
	CARBON_STEP       CarbonTransactionType = "STEP"       //步行
	CARBON_COFFEE_CUP CarbonTransactionType = "COFFEE_CUP" //自带咖啡杯
	CARBON_BIKE_RIDE  CarbonTransactionType = "BIKE_RIDE"  //骑行
	CARBON_ECAR       CarbonTransactionType = "ECAR"       //电动车主
)

type CarbonTransaction struct {
	ID            int64                 `json:"id"`
	OpenId        string                `gorm:"column:openid" json:"openId"`
	UserId        int64                 `json:"userId"`
	TransactionId string                `json:"transactionId"`
	Type          CarbonTransactionType `json:"type"`
	City          string                `json:"city"`
	Value         float64               `json:"value"`
	Info          string                `json:"info"`
	AdminId       int                   `json:"adminId"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
}
