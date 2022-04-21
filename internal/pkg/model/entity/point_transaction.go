package entity

import (
	"mio/internal/pkg/model"
	duibaApi "mio/pkg/duiba/api"
)

type PointTransactionType string

// Text 获取积分类型的中文名称
func (p PointTransactionType) Text() string {
	switch p {
	case POINT_STEP:
		return "步行"
	case POINT_COFFEE_CUP:
		return "自带咖啡杯"
	case POINT_PURCHASE:
		return "兑换"
	case POINT_INVITE:
		return "邀请好友"
	case POINT_CHECK_IN:
		return "每日签到"
	case POINT_BIKE_RIDE:
		return "骑行"
	case POINT_ECAR:
		return "电动车主"
	case POINT_COUPON:
		return "券码兑换"
	case POINT_QUIZ:
		return "答题活动"
	case POINT_PARTNERSHIP:
		return "合作活动"
	case POINT_GREEN_TORCH:
		return "绿炬人抽奖"
	case POINT_ADJUSTMENT:
		return "积分调整"
	case POINT_DUIBA_HDTOOL:
		return "活动抽奖"
	case POINT_DUIBA_GAME:
		return "游戏"
	case POINT_DUIBA_SIGN:
		return "活动奖励"
	case POINT_DUIBA_ALIPAY, POINT_DUIBA_QB, POINT_DUIBA_COUPON, POINT_DUIBA_OBJECT, POINT_DUIBA_PHONEBILL, POINT_DUIBA_PHONEFLOW, POINT_DUIBA_VIRTUAL:
		return "积分兑换"
	case POINT_DUIBA_REFUND:
		return "积分退还"
	case POINT_SYSTEM_REDUCE:
		return "系统扣减"
	case POINT_SYSTEM_ADD:
		return "系统补发"
	}
	return "未知积分"
}

const (
	POINT_STEP            PointTransactionType = "STEP"            //步行
	POINT_COFFEE_CUP      PointTransactionType = "COFFEE_CUP"      //自带咖啡杯
	POINT_PURCHASE        PointTransactionType = "PURCHASE"        //兑换
	POINT_INVITE          PointTransactionType = "INVITE"          //邀请好友
	POINT_CHECK_IN        PointTransactionType = "CHECK_IN"        //每日签到
	POINT_BIKE_RIDE       PointTransactionType = "BIKE_RIDE"       //骑行
	POINT_ECAR            PointTransactionType = "ECAR"            //电动车主
	POINT_COUPON          PointTransactionType = "COUPON"          //券码兑换
	POINT_QUIZ            PointTransactionType = "QUIZ"            //答题活动
	POINT_PARTNERSHIP     PointTransactionType = "PARTNERSHIP"     //合作活动
	POINT_GREEN_TORCH     PointTransactionType = "GREEN_TORCH"     //绿炬人抽奖
	POINT_ADJUSTMENT      PointTransactionType = "ADJUSTMENT"      //积分调整
	POINT_DUIBA_ALIPAY    PointTransactionType = "DUIBA_ALIPAY"    //积分兑换
	POINT_DUIBA_QB        PointTransactionType = "DUIBA_QB"        //积分兑换
	POINT_DUIBA_COUPON    PointTransactionType = "DUIBA_COUPON"    //积分兑换
	POINT_DUIBA_OBJECT    PointTransactionType = "DUIBA_OBJECT"    //积分兑换
	POINT_DUIBA_PHONEBILL PointTransactionType = "DUIBA_PHONEBILL" //积分兑换
	POINT_DUIBA_PHONEFLOW PointTransactionType = "DUIBA_PHONEFLOW" //积分兑换
	POINT_DUIBA_VIRTUAL   PointTransactionType = "DUIBA_VIRTUAL"   //积分兑换
	POINT_DUIBA_GAME      PointTransactionType = "DUIBA_GAME"      //游戏
	POINT_DUIBA_HDTOOL    PointTransactionType = "DUIBA_HDTOOL"    //活动抽奖
	POINT_DUIBA_SIGN      PointTransactionType = "DUIBA_SIGN"      //活动奖励
	POINT_DUIBA_REFUND    PointTransactionType = "DUIBA_REFUND"    //积分退还
	POINT_SYSTEM_REDUCE   PointTransactionType = "SYSTEM_REDUCE"   //系统扣减
	POINT_SYSTEM_ADD      PointTransactionType = "SYSTEM_ADD"      //系统补发
)

var PointTransactionTypeList = []PointTransactionType{
	POINT_STEP,
	POINT_COFFEE_CUP,
	POINT_PURCHASE,
	POINT_INVITE,
	POINT_CHECK_IN,
	POINT_BIKE_RIDE,
	POINT_ECAR,
	POINT_COUPON,
	POINT_QUIZ,
	POINT_PARTNERSHIP,
	POINT_GREEN_TORCH,
	POINT_ADJUSTMENT,
	POINT_DUIBA_ALIPAY,
	POINT_DUIBA_QB,
	POINT_DUIBA_COUPON,
	POINT_DUIBA_OBJECT,
	POINT_DUIBA_PHONEBILL,
	POINT_DUIBA_PHONEFLOW,
	POINT_DUIBA_VIRTUAL,
	POINT_DUIBA_GAME,
	POINT_DUIBA_HDTOOL,
	POINT_DUIBA_SIGN,
	POINT_DUIBA_REFUND,
	POINT_SYSTEM_REDUCE,
	POINT_SYSTEM_ADD,
}

const (
	OrderByPointTranCTDESC OrderBy = "order_by_point_ct_desc"
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
	ID             int64                `json:"id"`
	OpenId         string               `gorm:"column:openid" json:"openId"`
	TransactionId  string               `json:"transactionId"`
	Type           PointTransactionType `json:"type"`
	Value          int                  `json:"value"`
	CreateTime     model.Time           `json:"createTime"`
	AdditionalInfo AdditionalInfo       `json:"additionalInfo"`
	AdminId        int                  `json:"adminId"`
	Note           string               `json:"note"`
}
type AdditionalInfo string

type PointPurchaseInfo struct {
}

func (info AdditionalInfo) ToDuiBa() duibaApi.ExchangeForm {
	return duibaApi.ExchangeForm{}
}
func (info AdditionalInfo) ToDuiBaRefund() duibaApi.ExchangeResultForm {
	return duibaApi.ExchangeResultForm{}
}
func (info AdditionalInfo) ToPurchase() PointPurchaseInfo {
	return PointPurchaseInfo{}
}
