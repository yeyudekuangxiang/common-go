package entity

import "mio/model"

type PointTransactionType string

const (
	POINT_STEP        PointTransactionType = "STEP"        //步行
	POINT_COFFEE_CUP  PointTransactionType = "COFFEE_CUP"  //自带咖啡杯
	POINT_PURCHASE    PointTransactionType = "PURCHASE"    //兑换
	POINT_INVITE      PointTransactionType = "INVITE"      //邀请好友
	POINT_CHECK_IN    PointTransactionType = "CHECK_IN"    //每日签到
	POINT_BIKE_RIDE   PointTransactionType = "BIKE_RIDE"   //骑行
	POINT_ECAR        PointTransactionType = "ECAR"        //电动车主
	POINT_COUPON      PointTransactionType = "COUPON"      //券码兑换
	POINT_QUIZ        PointTransactionType = "QUIZ"        //答题活动
	POINT_PARTNERSHIP PointTransactionType = "PARTNERSHIP" //合作活动
	POINT_GREEN_TORCH PointTransactionType = "GREEN_TORCH" //绿炬人抽奖
	POINT_ADJUSTMENT  PointTransactionType = "ADJUSTMENT"  //积分调整
)

var PointCollectValueMap = map[PointTransactionType]int{
	POINT_COFFEE_CUP: 150,
	POINT_BIKE_RIDE:  150,
	POINT_ECAR:       300,
	POINT_INVITE:     500,
}
var PointCollectLimitMap = map[PointTransactionType]int{
	POINT_COFFEE_CUP: 4,
	POINT_BIKE_RIDE:  2,
	POINT_ECAR:       1,
	POINT_INVITE:     5,
}

type PointTransaction struct {
	Id             int64                `json:"id"`
	OpenId         string               `gorm:"column:openid" json:"openId"`
	TransactionId  string               `json:"transactionId"`
	Type           PointTransactionType `json:"type"`
	Value          int                  `json:"value"`
	CreateTime     model.Time           `json:"createTime"`
	AdditionalInfo string               `json:"additionalInfo"`
}
