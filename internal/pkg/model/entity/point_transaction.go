package entity

import (
	"database/sql"
	duibaApi "gitlab.miotech.com/miotech-application/backend/common-go/duiba/api/model"
	"mio/internal/pkg/model"
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
	case POINT_ECAR, POINT_FAST_ELECTRICITY:
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
	case POINT_PLATFORM:
		return "活动奖励" //平台活动
	case POINT_DUIBA_INTEGRAL_RECHARGE:
		return "活动奖励"
	case POINT_POWER_REPLACE:
		return "新能源换电" //oola

	case POINT_RECYCLING, POINT_RECYCLING_CLOTHING, POINT_RECYCLING_DIGITAL, POINT_RECYCLING_APPLIANCE, POINT_RECYCLING_BOOK, POINT_FMY_RECYCLING_CLOTHING:
		return "旧物回收"
	case POINT_ARTICLE:
		return "发帖"
	case POINT_RECOMMEND:
		return "推荐"
	case POINT_COMMENT:
		return "评论"
	case POINT_LIKE:
		return "点赞"
	case POINT_REDUCE_PLASTIC:
		return "环保减塑"
	case POINT_ZCYP_SIGNUP, POINT_ZCYP_APPLY:
		return "零碳小先锋"
	case POINT_CYCLING:
		return "骑行"
	case POINT_YTX:
		return "地铁出行"
	case POINT_SECONDHAND_ORDER:
		return "二手交易"
	case POINT_SECONDHAND_ORDER_AWARD:
		return "二手交易奖励"
	case POINT_NEW_TASK_PUBLISH_COMMODITY:
		return "发布二手奖励"
	case POINT_RECYCLING_SHISHANGHUISHOU:
		return "拾尚回收"
	case POINT_RECYCLING_DANGDANGYIXIA:
		return "铛铛一下"
	case POINT_HELLO_BIKE_RIDE:
		return "骑行"
	case POINT_ECAR_MIO:
		return "绿喵积分"
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
	case POINT_PLATFORM:
		return "平台活动奖励"
	case POINT_DUIBA_INTEGRAL_RECHARGE:
		return "兑吧积分充值"
	case POINT_RECYCLING_CLOTHING:
		return "oola衣物鞋帽回收"
	case POINT_RECYCLING_DIGITAL:
		return "oola数码回收"
	case POINT_RECYCLING_APPLIANCE:
		return "oola家电回收"
	case POINT_RECYCLING_BOOK:
		return "oola书籍回收"
	case POINT_ECAR:
		return "星星充电"
	case POINT_FAST_ELECTRICITY:
		return "快电"
	case POINT_CYCLING:
		return "骑行"
	case POINT_YTX:
		return "亿通行"
	case POINT_RECYCLING_AIHUISHOU:
		return "爱分类/爱回收"
	case POINT_RECYCLING:
		return "旧物回收"
	case POINT_RECYCLING_SHISHANGHUISHOU:
		return "拾尚回收"
	case POINT_RECYCLING_DANGDANGYIXIA:
		return "铛铛一下"
	case POINT_YKC:
		return "云快充"
	case POINT_HELLO_BIKE_RIDE:
		return "骑行"
	case POINT_ECAR_MIO:
		return "绿喵充电"
	}
	return p.Text()
}

const (
	POINT_STEP                    PointTransactionType = "STEP"              //步行
	POINT_COFFEE_CUP              PointTransactionType = "COFFEE_CUP"        //自带咖啡杯
	POINT_PURCHASE                PointTransactionType = "PURCHASE"          //兑换  所有走订单逻辑的都叫兑换：证书&券码兑换商品
	POINT_INVITE                  PointTransactionType = "INVITE"            //邀请好友
	POINT_CHECK_IN                PointTransactionType = "CHECK_IN"          //每日签到
	POINT_BIKE_RIDE               PointTransactionType = "BIKE_RIDE"         //骑行
	POINT_ECAR                    PointTransactionType = "ECAR"              //电动车主
	POINT_COUPON                  PointTransactionType = "COUPON"            //券码兑换
	POINT_QUIZ                    PointTransactionType = "QUIZ"              //答题活动
	POINT_PARTNERSHIP             PointTransactionType = "PARTNERSHIP"       //合作活动
	POINT_PLATFORM                PointTransactionType = "PLATFORM"          //合作活动
	POINT_GREEN_TORCH             PointTransactionType = "GREEN_TORCH"       //绿炬人抽奖
	POINT_ADJUSTMENT              PointTransactionType = "ADJUSTMENT"        //积分调整
	POINT_DUIBA_ALIPAY            PointTransactionType = "DUIBA_ALIPAY"      //兑吧支付宝 积分兑换
	POINT_DUIBA_QB                PointTransactionType = "DUIBA_QB"          //兑吧qb 积分兑换
	POINT_DUIBA_COUPON            PointTransactionType = "DUIBA_COUPON"      //兑吧优惠券 积分兑换
	POINT_DUIBA_OBJECT            PointTransactionType = "DUIBA_OBJECT"      //兑吧实物 积分兑换
	POINT_DUIBA_PHONEBILL         PointTransactionType = "DUIBA_PHONEBILL"   //兑吧话费 积分兑换
	POINT_DUIBA_PHONEFLOW         PointTransactionType = "DUIBA_PHONEFLOW"   //兑吧流量 积分兑换
	POINT_DUIBA_VIRTUAL           PointTransactionType = "DUIBA_VIRTUAL"     //兑吧虚拟商品 积分兑换
	POINT_DUIBA_GAME              PointTransactionType = "DUIBA_GAME"        //兑吧游戏 游戏
	POINT_DUIBA_HDTOOL            PointTransactionType = "DUIBA_HDTOOL"      //兑吧活动抽奖 活动抽奖
	POINT_DUIBA_SIGN              PointTransactionType = "DUIBA_SIGN"        //兑吧签到 活动奖励
	POINT_DUIBA_REFUND            PointTransactionType = "DUIBA_REFUND"      //兑吧积分退还 积分退还
	POINT_DUIBA_POSTSALE          PointTransactionType = "DUIBA_POSTSALE"    //售后退积分 退积分
	POINT_DUIBA_CANCELSHIP        PointTransactionType = "DUIBA_CANCELSHIP"  //取消发货 退积分
	POINT_DUIBA_TASK              PointTransactionType = "DUIBA_TASK"        //pk比赛 pk比赛
	POINT_SYSTEM_REDUCE           PointTransactionType = "SYSTEM_REDUCE"     //系统扣减
	POINT_SYSTEM_ADD              PointTransactionType = "SYSTEM_ADD"        //系统补发
	POINT_JHX                     PointTransactionType = "JHX"               //金华行
	POINT_POWER_REPLACE           PointTransactionType = "POWER_REPLACE"     //电车换电
	POINT_DUIBA_INTEGRAL_RECHARGE PointTransactionType = "INTEGRAL_RECHARGE" //兑吧虚拟商品充值积分
	POINT_YTX                     PointTransactionType = "YTX"               //亿通行
	POINT_YKC                     PointTransactionType = "YKC"

	POINT_RECYCLING                 PointTransactionType = "RECYCLING"                 //旧物回收
	POINT_RECYCLING_AIHUISHOU       PointTransactionType = "RECYCLING_AIHUISHOU"       //旧物回收 爱分类爱回收
	POINT_RECYCLING_CLOTHING        PointTransactionType = "RECYCLING_CLOTHING"        //旧物回收 oola衣物鞋帽
	POINT_RECYCLING_DIGITAL         PointTransactionType = "RECYCLING_COMPUTER"        //旧物回收 oola数码
	POINT_RECYCLING_APPLIANCE       PointTransactionType = "RECYCLING_APPLIANCE"       //旧物回收 oola家电
	POINT_RECYCLING_BOOK            PointTransactionType = "RECYCLING_BOOK"            //旧物回收 oola书籍
	POINT_FMY_RECYCLING_CLOTHING    PointTransactionType = "RECYCLING_FMY_CLOTHING"    //旧物回收 fmy衣物鞋帽
	POINT_RECYCLING_SHISHANGHUISHOU PointTransactionType = "RECYCLING_SHISHANGHUISHOU" //旧物回收 拾尚回收
	POINT_RECYCLING_DANGDANGYIXIA   PointTransactionType = "RECYCLING_DANGDANGYIXIA"   //旧物回收 铛铛一下

	POINT_FAST_ELECTRICITY PointTransactionType = "FAST_ELECTRICITY" //快电

	POINT_ARTICLE   PointTransactionType = "ARTICLE"   //发文章
	POINT_COMMENT   PointTransactionType = "COMMENT"   //评论
	POINT_RECOMMEND PointTransactionType = "RECOMMEND" //推荐
	POINT_LIKE      PointTransactionType = "LIKE"      //点赞

	POINT_REDUCE_PLASTIC PointTransactionType = "REDUCE_PLASTIC" //环保减塑

	POINT_ZCYP_SIGNUP PointTransactionType = "ZCYP_SIGNUP"
	POINT_ZCYP_APPLY  PointTransactionType = "ZCYP_APPLY"

	POINT_CYCLING                    PointTransactionType = "CYCLING"
	POINT_SECONDHAND_ORDER           PointTransactionType = "SECONDHAND_ORDER"
	POINT_SECONDHAND_ORDER_AWARD     PointTransactionType = "SECONDHAND_ORDER_AWARD"
	POINT_NEW_TASK_PUBLISH_COMMODITY PointTransactionType = "NEW_TASK_PUBLISH_COMMODITY"
	POINT_HELLO_BIKE_RIDE            PointTransactionType = "HELLO_BIKE_RIDE" //哈啰骑行
	POINT_ECAR_MIO                   PointTransactionType = "ECAR_MIO"
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
	POINT_PLATFORM,
	POINT_DUIBA_POSTSALE,
	POINT_DUIBA_CANCELSHIP,
	POINT_DUIBA_TASK,
	POINT_ARTICLE,
	POINT_RECOMMEND,
	POINT_DUIBA_INTEGRAL_RECHARGE,
	POINT_RECYCLING_CLOTHING,
	POINT_RECYCLING_DIGITAL,
	POINT_RECYCLING_APPLIANCE,
	POINT_RECYCLING_BOOK,
	POINT_FAST_ELECTRICITY,
	POINT_COMMENT,
	POINT_LIKE,
	POINT_REDUCE_PLASTIC,
	POINT_CYCLING,
	POINT_SECONDHAND_ORDER,
	POINT_SECONDHAND_ORDER_AWARD,
	POINT_NEW_TASK_PUBLISH_COMMODITY,
	POINT_HELLO_BIKE_RIDE,
	POINT_ECAR_MIO,
}

const (
	OrderByPointTranCTDESC OrderBy = "order_by_point_ct_desc"
)

// 积分限制
var PointCollectValueMap = map[PointTransactionType]int{
	POINT_COFFEE_CUP:     39,  //	每次
	POINT_BIKE_RIDE:      42,  //	每次
	POINT_INVITE:         500, //	每人/天
	POINT_POWER_REPLACE:  300, //	每人/天
	POINT_ARTICLE:        150, //	每次
	POINT_COMMENT:        10,  //每次
	POINT_LIKE:           5,   //每次
	POINT_RECOMMEND:      50,  //每次
	POINT_REDUCE_PLASTIC: 28,  //每次
	POINT_ZCYP_SIGNUP:    2500,
	POINT_ZCYP_APPLY:     50,
}

var PointCollectLimitOnceMap = map[PointTransactionType]int{
	POINT_ZCYP_SIGNUP: 1,
	POINT_ZCYP_APPLY:  1,
}

var PointTypesMap = map[string]PointTransactionType{
	"yitongxing":  POINT_YTX,
	"jinghuaxing": POINT_JHX,
}

//每天获取 （多少）次积分
var PointCollectLimitMap = map[PointTransactionType]int{
	POINT_COFFEE_CUP:    2, //	次
	POINT_BIKE_RIDE:     2, //	次
	POINT_INVITE:        5, //	次
	POINT_POWER_REPLACE: 1, //	次
	//POINT_ARTICLE:        2, //	次
	POINT_COMMENT:        3, // 次
	POINT_LIKE:           6, // 次
	POINT_REDUCE_PLASTIC: 2, // 次
}

//b端渠道对应操作类型
var PlatformMethodMap = map[string]PointTransactionType{
	"zcyp_signup": POINT_ZCYP_SIGNUP,
	"zcyp_apply":  POINT_ZCYP_APPLY,
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
	ExpireTime     sql.NullTime         `json:"expireTime"`
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
