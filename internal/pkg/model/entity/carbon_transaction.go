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
	case CARBON_TESR:
		return "测试"
	}
	return "未知积分"
}

func (p CarbonTransactionType) Cover() string {
	switch p {
	case CARBON_STEP:
		return "step_url"
	case CARBON_COFFEE_CUP:
		return "cup_url"
	case CARBON_BIKE_RIDE:
		return "bike_url"
	case CARBON_ECAR:
		return "car_url"
	case CARBON_TESR:
		return "test_url"
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
	case CARBON_TESR:
		return "测试"
	}
	return p.Text()
}

const (
	CARBON_STEP       CarbonTransactionType = "STEP"       //步行
	CARBON_COFFEE_CUP CarbonTransactionType = "COFFEE_CUP" //自带咖啡杯
	CARBON_BIKE_RIDE  CarbonTransactionType = "BIKE_RIDE"  //骑行
	CARBON_ECAR       CarbonTransactionType = "ECAR"       //电动车主
	CARBON_TESR       CarbonTransactionType = "TEST"       //测试用
)

var CarbonTransactionTypeList = []CarbonTransactionType{
	CARBON_STEP,
	CARBON_COFFEE_CUP,
}

var CarbonCollectValueMap = map[CarbonTransactionType]int{
	CARBON_STEP:       150,
	CARBON_COFFEE_CUP: 150,
}

var CarbonCollectLimitMap = map[CarbonTransactionType]int{
	CARBON_STEP:       4,
	CARBON_COFFEE_CUP: 2,
}

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
