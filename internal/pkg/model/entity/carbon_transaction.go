package entity

import (
	"mio/config"
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
	case CARBON_POWER_REPLACE:
		return "换电"
	case CARBON_REDUCE_PLASTIC:
		return "环保减塑"
	case CARBON_RECYCLING_CLOTHING:
		return "旧物回收 oola衣物鞋帽"
	case CARBON_RECYCLING_DIGITAL:
		return "旧物回收 oola数码"
	case CARBON_RECYCLING_APPLIANCE:
		return "旧物回收 oola家电"
	case CARBON_RECYCLING_BOOK:
		return "旧物回收 oola书籍"
	case CARBON_FMY_RECYCLING_CLOTHING:
		return "旧物回收 fmy衣物鞋帽"
	case CARBON_RECYCLING:
		return "旧物回收"
	case CARBON_JHX:
		return "金华行"
	case CARBON_YTX:
		return "地铁"

	}
	return "未知积分"
}

func (p CarbonTransactionType) Cover() string {
	switch p {
	case CARBON_STEP:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_foot.png"
	case CARBON_COFFEE_CUP:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_cup.png"
	case CARBON_BIKE_RIDE:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_riding.png"
	case CARBON_ECAR:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_newenergy.png"
	case CARBON_YTX:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_metro.png"
	case CARBON_JHX:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_bus.png"
	case CARBON_RECYCLING_CLOTHING:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_declutter.png"
	case CARBON_RECYCLING_DIGITAL:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_declutter.png"
	case CARBON_RECYCLING_APPLIANCE:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_declutter.png"
	case CARBON_RECYCLING_BOOK:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_declutter.png"
	case CARBON_FMY_RECYCLING_CLOTHING:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_declutter.png"
	case CARBON_RECYCLING:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_declutter.png"
	case CARBON_POWER_REPLACE:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_newenergy.png"
	case CARBON_REDUCE_PLASTIC:
		return config.Config.OSS.CdnDomain + "/static/mp2c/images/mp2c2.0/assets/classify_newenergy.png"
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
	case CARBON_POWER_REPLACE:
		return "换电"
	case CARBON_REDUCE_PLASTIC:
		return "环保减塑"
	case CARBON_RECYCLING_CLOTHING:
		return "旧物回收 oola衣物鞋帽"
	case CARBON_RECYCLING_DIGITAL:
		return "旧物回收 oola数码"
	case CARBON_RECYCLING_APPLIANCE:
		return "旧物回收 oola家电"
	case CARBON_RECYCLING_BOOK:
		return "旧物回收 oola书籍"
	case CARBON_FMY_RECYCLING_CLOTHING:
		return "旧物回收 fmy衣物鞋帽"
	case CARBON_JHX:
		return "金华行"
	case CARBON_RECYCLING:
		return "旧物回收"
	case CARBON_CYCLING:
		return "骑行"
	case CARBON_YTX:
		return "地铁"
	case CARBON_RECYCLING_AIHUISHOU:
		return "爱分类/爱回收"
	case CARBON_RECYCLING_SHISHANGHUISHOU:
		return "拾尚回收"
	case CARBON_RECYCLING_DANGDANGYIXIA:
		return "铛铛一下"
	case CARBON_YKC:
		return "云快充"
	}
	return p.Text()
}

const (
	CARBON_STEP                      CarbonTransactionType = "STEP"                   //步行
	CARBON_COFFEE_CUP                CarbonTransactionType = "COFFEE_CUP"             //自带咖啡杯
	CARBON_BIKE_RIDE                 CarbonTransactionType = "BIKE_RIDE"              //骑行
	CARBON_ECAR                      CarbonTransactionType = "ECAR"                   //电动车主
	CARBON_POWER_REPLACE             CarbonTransactionType = "POWER_REPLACE"          //换电
	CARBON_REDUCE_PLASTIC            CarbonTransactionType = "REDUCE_PLASTIC"         //环保减塑
	CARBON_RECYCLING_CLOTHING        CarbonTransactionType = "RECYCLING_CLOTHING"     //旧物回收 oola衣物鞋帽
	CARBON_RECYCLING_DIGITAL         CarbonTransactionType = "RECYCLING_COMPUTER"     //旧物回收 oola数码
	CARBON_RECYCLING_APPLIANCE       CarbonTransactionType = "RECYCLING_APPLIANCE"    //旧物回收 oola家电
	CARBON_RECYCLING_BOOK            CarbonTransactionType = "RECYCLING_BOOK"         //旧物回收 oola书籍
	CARBON_FMY_RECYCLING_CLOTHING    CarbonTransactionType = "RECYCLING_FMY_CLOTHING" //旧物回收 fmy衣物鞋帽
	CARBON_JHX                       CarbonTransactionType = "JHX"                    //金华行
	CARBON_RECYCLING                 CarbonTransactionType = "RECYCLING"              //旧物回收总的
	CARBON_RECYCLING_AIHUISHOU       CarbonTransactionType = "RECYCLING_AIHUISHOU"    //爱回收
	CARBON_RECYCLING_SHISHANGHUISHOU CarbonTransactionType = "RECYCLING_SHISHANGHUISHOU"
	CARBON_RECYCLING_DANGDANGYIXIA   CarbonTransactionType = "RECYCLING_DANGDANGYIXIA"
	CARBON_YTX                       CarbonTransactionType = "YTX"     //亿通行
	CARBON_CYCLING                   CarbonTransactionType = "CYCLING" //java 骑行
	CARBON_YKC                       CarbonTransactionType = "YKC"     //云快充

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
