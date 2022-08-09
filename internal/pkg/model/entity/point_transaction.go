package entity

import (
	"mio/internal/pkg/model"
	duibaApi "mio/pkg/duiba/api/model"
)

type PointTransactionType string

// Text 获取积分类型的中文名称(给用户看的)
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
		return "新能源充电"
	case POINT_COUPON:
		return "券码兑换"
	case POINT_QUIZ:
		return "答题"
	case POINT_PARTNERSHIP:
		return "活动奖励"
	case POINT_GREEN_TORCH:
		return "绿炬人抽奖"
	case POINT_ADJUSTMENT:
		return "积分调整"
	case POINT_DUIBA_HDTOOL:
		return "活动抽奖"
	case POINT_DUIBA_GAME:
		return "游戏奖励"
	case POINT_DUIBA_SIGN:
		return "签到"
	case POINT_DUIBA_ALIPAY, POINT_DUIBA_QB, POINT_DUIBA_COUPON, POINT_DUIBA_OBJECT, POINT_DUIBA_PHONEBILL, POINT_DUIBA_PHONEFLOW, POINT_DUIBA_VIRTUAL:
		return "积分兑换"
	case POINT_DUIBA_REFUND:
		return "退积分"
	case POINT_DUIBA_TASK:
		return "PK赛"
	case POINT_DUIBA_CANCELSHIP:
		return "取消发货"
	case POINT_DUIBA_POSTSALE:
		return "退积分"
	case POINT_SYSTEM_REDUCE:
		return "系统扣减"
	case POINT_SYSTEM_ADD:
		return "系统补发"
	case POINT_JHX:
		return "公交出行"
	}
	return "未知积分"
}

// RealText 获取积分类型的中文名称(给管理员看的)
func (p PointTransactionType) RealText() string {
	switch p {
	case POINT_DUIBA_ALIPAY:
		return "兑吧支付宝"
	case POINT_DUIBA_QB:
		return "兑吧qb"
	case POINT_DUIBA_COUPON:
		return "兑吧优惠券"
	case POINT_DUIBA_OBJECT:
		return "兑吧实物"
	case POINT_DUIBA_PHONEBILL:
		return "兑吧话费"
	case POINT_DUIBA_PHONEFLOW:
		return "兑吧流量"
	case POINT_DUIBA_VIRTUAL:
		return "兑吧虚拟商品"
	case POINT_DUIBA_REFUND:
		return "兑吧积分退还"
	case POINT_JHX:
		return "金华行"
	}
	return p.Text()
}

const (
	POINT_STEP            PointTransactionType = "STEP"            //步行
	POINT_COFFEE_CUP      PointTransactionType = "COFFEE_CUP"      //自带咖啡杯
	POINT_PURCHASE        PointTransactionType = "PURCHASE"        //兑换
	POINT_INVITE          PointTransactionType = "INVITE"          //邀请好友
	POINT_CHECK_IN        PointTransactionType = "CHECK_IN"        //每日签到
	POINT_BIKE_RIDE       PointTransactionType = "BIKE_RIDE"       //骑行
	POINT_ECAR            PointTransactionType = "ECAR"            //电动车主
	POINT_POWER_REPLACE   PointTransactionType = "POWER_REPLACE"   //电车换电
	POINT_COUPON          PointTransactionType = "COUPON"          //券码兑换
	POINT_QUIZ            PointTransactionType = "QUIZ"            //答题活动
	POINT_PARTNERSHIP     PointTransactionType = "PARTNERSHIP"     //合作活动
	POINT_GREEN_TORCH     PointTransactionType = "GREEN_TORCH"     //绿炬人抽奖
	POINT_ADJUSTMENT      PointTransactionType = "ADJUSTMENT"      //积分调整
	POINT_DUIBA_ALIPAY    PointTransactionType = "DUIBA_ALIPAY"    //兑吧支付宝 积分兑换
	POINT_DUIBA_QB        PointTransactionType = "DUIBA_QB"        //兑吧qb 积分兑换
	POINT_DUIBA_COUPON    PointTransactionType = "DUIBA_COUPON"    //兑吧优惠券 积分兑换
	POINT_DUIBA_OBJECT    PointTransactionType = "DUIBA_OBJECT"    //兑吧实物 积分兑换
	POINT_DUIBA_PHONEBILL PointTransactionType = "DUIBA_PHONEBILL" //兑吧话费 积分兑换
	POINT_DUIBA_PHONEFLOW PointTransactionType = "DUIBA_PHONEFLOW" //兑吧流量 积分兑换
	POINT_DUIBA_VIRTUAL   PointTransactionType = "DUIBA_VIRTUAL"   //兑吧虚拟商品 积分兑换
	POINT_DUIBA_GAME      PointTransactionType = "DUIBA_GAME"      //兑吧游戏 游戏
	POINT_DUIBA_HDTOOL    PointTransactionType = "DUIBA_HDTOOL"    //兑吧活动抽奖 活动抽奖
	POINT_DUIBA_SIGN      PointTransactionType = "DUIBA_SIGN"      //兑吧签到 活动奖励
	POINT_DUIBA_REFUND    PointTransactionType = "DUIBA_REFUND"    //兑吧积分退还 积分退还
	POINT_SYSTEM_REDUCE   PointTransactionType = "SYSTEM_REDUCE"   //系统扣减
	POINT_SYSTEM_ADD      PointTransactionType = "SYSTEM_ADD"      //系统补发
	POINT_JHX             PointTransactionType = "JHX"             //金华行
	POINT_ARTICAL         PointTransactionType = "ARTICAL"         //发文章
	POINT_RECOMMEND       PointTransactionType = "RECOMMEND"       //文章/评论被推荐
	POINT_DUIBA_POSTSALE   PointTransactionType = "DUIBA_POSTSALE"   //售后退积分 退积分
	POINT_DUIBA_CANCELSHIP PointTransactionType = "DUIBA_CANCELSHIP" //取消发货 退积分
	POINT_DUIBA_TASK       PointTransactionType = "DUIBA_TASK"       //pk比赛 pk比赛
)

var PointTransactionTypeList = []PointTransactionType{
	POINT_STEP,
	POINT_COFFEE_CUP,
	POINT_POWER_REPLACE,
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
	POINT_DUIBA_POSTSALE,
	POINT_DUIBA_CANCELSHIP,
	POINT_DUIBA_TASK,
	POINT_ARTICAL,
	POINT_RECOMMEND,
}

const (
	OrderByPointTranCTDESC OrderBy = "order_by_point_ct_desc"
)

//每次获取 （多少） 积分
var PointCollectValueMap = map[PointTransactionType]int{
	POINT_COFFEE_CUP:    39,  //每次
	POINT_BIKE_RIDE:     42,  //每次
	POINT_INVITE:        500, //每人
	POINT_POWER_REPLACE: 300, //每次
	POINT_ARTICAL:       150, //每次
}

//每天获取 （多少）次积分
var PointCollectLimitMap = map[PointTransactionType]int{
	POINT_COFFEE_CUP:    2,
	POINT_BIKE_RIDE:     2,
	POINT_INVITE:        5,
	POINT_POWER_REPLACE: 1,
	POINT_ARTICAL:       2,
}

type PointTransaction struct {
	ID             int64                `json:"id"`
	OpenId         string               `gorm:"column:openid" json:"openId"`
	TransactionId  string               `json:"transactionId"`
	Type           PointTransactionType `json:"type"`
	Value          int64                `json:"value"`
	CreateTime     model.Time           `json:"createTime"`
	AdditionalInfo AdditionalInfo       `json:"additionalInfo"`
	AdminId        int                  `json:"adminId"`
	Note           string               `json:"note"`
}
type AdditionalInfo string

type PointPurchaseInfo struct {
}

func (info AdditionalInfo) ToDuiBa() duibaApi.Exchange {
	return duibaApi.Exchange{}
}
func (info AdditionalInfo) ToDuiBaRefund() duibaApi.ExchangeResult {
	return duibaApi.ExchangeResult{}
}
func (info AdditionalInfo) ToPurchase() PointPurchaseInfo {
	return PointPurchaseInfo{}
}
