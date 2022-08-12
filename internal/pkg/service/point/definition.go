package point

type CommandDescription struct {
	Times  int64                     //限制次数
	Amount int64                     //每次数量
	Fn     func(*clientHandle) error //执行方法
}

type CollectType string
type CollectRules map[CollectType][]string

//type CollectError map[CollectType]error

const (
	STEP             CollectType = "STEP"             //步行
	COFFEE_CUP       CollectType = "COFFEE_CUP"       //自带咖啡杯
	PURCHASE         CollectType = "PURCHASE"         //兑换
	INVITE           CollectType = "INVITE"           //邀请好友
	CHECK_IN         CollectType = "CHECK_IN"         //每日签到
	BIKE_RIDE        CollectType = "BIKE_RIDE"        //骑行
	ECAR             CollectType = "ECAR"             //电动车主
	POWER_REPLACE    CollectType = "POWER_REPLACE"    //电车换电
	COUPON           CollectType = "COUPON"           //券码兑换
	QUIZ             CollectType = "QUIZ"             //答题活动
	PARTNERSHIP      CollectType = "PARTNERSHIP"      //合作活动
	GREEN_TORCH      CollectType = "GREEN_TORCH"      //绿炬人抽奖
	ADJUSTMENT       CollectType = "ADJUSTMENT"       //积分调整
	DUIBA_ALIPAY     CollectType = "DUIBA_ALIPAY"     //兑吧支付宝 积分兑换
	DUIBA_QB         CollectType = "DUIBA_QB"         //兑吧qb 积分兑换
	DUIBA_COUPON     CollectType = "DUIBA_COUPON"     //兑吧优惠券 积分兑换
	DUIBA_OBJECT     CollectType = "DUIBA_OBJECT"     //兑吧实物 积分兑换
	DUIBA_PHONEBILL  CollectType = "DUIBA_PHONEBILL"  //兑吧话费 积分兑换
	DUIBA_PHONEFLOW  CollectType = "DUIBA_PHONEFLOW"  //兑吧流量 积分兑换
	DUIBA_VIRTUAL    CollectType = "DUIBA_VIRTUAL"    //兑吧虚拟商品 积分兑换
	DUIBA_GAME       CollectType = "DUIBA_GAME"       //兑吧游戏 游戏
	DUIBA_HDTOOL     CollectType = "DUIBA_HDTOOL"     //兑吧活动抽奖 活动抽奖
	DUIBA_SIGN       CollectType = "DUIBA_SIGN"       //兑吧签到 活动奖励
	DUIBA_REFUND     CollectType = "DUIBA_REFUND"     //兑吧积分退还 积分退还
	SYSTEM_REDUCE    CollectType = "SYSTEM_REDUCE"    //系统扣减
	SYSTEM_ADD       CollectType = "SYSTEM_ADD"       //系统补发
	JHX              CollectType = "JHX"              //金华行
	ARTICLE          CollectType = "ARTICLE"          //发文章
	RECOMMEND        CollectType = "RECOMMEND"        //文章/评论被推荐
	DUIBA_POSTSALE   CollectType = "DUIBA_POSTSALE"   //售后退积分 退积分
	DUIBA_CANCELSHIP CollectType = "DUIBA_CANCELSHIP" //取消发货 退积分
	DUIBA_TASK       CollectType = "DUIBA_TASK"       //pk比赛 pk比赛
	PLATFORM         CollectType = "PLATFORM"         //活动奖励
)

var rules = CollectRules{
	"COFFEE_CUP":    []string{"自带杯", "单号", "订单"},
	"BIKE_RIDE":     []string{"骑行", "单车", "骑车", "bike", "出行", "哈啰", "摩拜", "青桔"},
	"POWER_REPLACE": []string{"订单编号", "已支付"},
}

//var collectError = CollectError{
//	"COFFEE_CUP":    errors.New(""),
//	"BIKE_RIDE":     errors.New(""),
//	"POWER_REPLACE": errors.New(""),
//}

var commandText = map[CollectType]string{
	DUIBA_ALIPAY:     "积分兑换",
	DUIBA_QB:         "积分兑换",
	DUIBA_COUPON:     "积分兑换",
	DUIBA_OBJECT:     "积分兑换",
	DUIBA_PHONEBILL:  "积分兑换",
	DUIBA_PHONEFLOW:  "积分兑换",
	DUIBA_VIRTUAL:    "积分兑换",
	STEP:             "步行",
	COFFEE_CUP:       "自带咖啡杯",
	PURCHASE:         "兑换",
	INVITE:           "邀请好友",
	CHECK_IN:         "每日签到",
	BIKE_RIDE:        "骑行",
	ECAR:             "新能源充电",
	COUPON:           "券码兑换",
	QUIZ:             "答题",
	PARTNERSHIP:      "活动奖励",
	GREEN_TORCH:      "绿炬人抽奖",
	ADJUSTMENT:       "积分调整",
	DUIBA_HDTOOL:     "活动抽奖",
	DUIBA_GAME:       "游戏奖励",
	DUIBA_SIGN:       "签到",
	DUIBA_REFUND:     "退积分",
	DUIBA_TASK:       "PK赛",
	DUIBA_CANCELSHIP: "取消发货",
	DUIBA_POSTSALE:   "退积分",
	SYSTEM_REDUCE:    "系统扣减",
	SYSTEM_ADD:       "系统补发",
	JHX:              "公交出行",
	PLATFORM:         "活动奖励",
	ARTICLE:          "文章",
	RECOMMEND:        "推荐",
	POWER_REPLACE:    "电车换电",
}

var commandRealText = map[CollectType]string{
	DUIBA_ALIPAY:    "兑吧支付宝",
	DUIBA_QB:        "兑吧qb",
	DUIBA_COUPON:    "兑吧优惠券",
	DUIBA_OBJECT:    "兑吧实物",
	DUIBA_PHONEBILL: "兑吧话费",
	DUIBA_PHONEFLOW: "兑吧流量",
	DUIBA_VIRTUAL:   "兑吧虚拟商品",
	DUIBA_REFUND:    "兑吧积分退还",
	JHX:             "金华行",
	PLATFORM:        "平台活动奖励",
}

//方法定义
var commandMap = map[string]*CommandDescription{
	"COFFEE_CUP":       {Fn: (*clientHandle).coffeeCup, Times: 2, Amount: 39},
	"INVITE":           {Fn: (*clientHandle).invite, Times: 5, Amount: 500},
	"BIKE_RIDE":        {Fn: (*clientHandle).bikeRide, Times: 2, Amount: 42},
	"POWER_REPLACE":    {Fn: (*clientHandle).powerReplace, Times: 1, Amount: 300},
	"ARTICLE":          {Fn: (*clientHandle).article, Times: 2, Amount: 150},
	"FAST_ELECTRICITY": {Fn: (*clientHandle).fastElectricity, Times: 0, Amount: 300},
}

//dui ba
//var depleteCommandMap = map[string]func(){}
