package point

type CommandDescription struct {
	Number int                                       //当前次数
	Open   bool                                      // 是否实现
	Fn     func(*clientHandle, ...interface{}) error //执行方法
}

type CollectType string
type CollectRules map[CollectType][]string

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
	ARTICAL          CollectType = "ARTICAL"          //发文章
	RECOMMEND        CollectType = "RECOMMEND"        //文章/评论被推荐
	DUIBA_POSTSALE   CollectType = "DUIBA_POSTSALE"   //售后退积分 退积分
	DUIBA_CANCELSHIP CollectType = "DUIBA_CANCELSHIP" //取消发货 退积分
	DUIBA_TASK       CollectType = "DUIBA_TASK"       //pk比赛 pk比赛
)

var rules = CollectRules{
	"COFFEE_CUP":    []string{"自带杯", "单号", "订单"},
	"BIKE_RIDE":     []string{"骑行", "单车", "骑车", "bike", "出行", "哈啰", "摩拜", "青桔"},
	"POWER_REPLACE": []string{"订单编号"},
}

var commandMap = map[string]*CommandDescription{
	"STEP":             {},
	"COFFEE_CUP":       {Fn: (*clientHandle).coffeeCup, Open: true},
	"PURCHASE":         {},
	"INVITE":           {},
	"CHECK_IN":         {},
	"BIKE_RIDE":        {},
	"ECAR":             {},
	"POWER_REPLACE":    {Fn: (*clientHandle).powerReplace, Open: true},
	"COUPON":           {},
	"QUIZ":             {},
	"PARTNERSHIP":      {},
	"GREEN_TORCH":      {},
	"ADJUSTMENT":       {},
	"DUIBA_ALIPAY":     {},
	"DUIBA_QB":         {},
	"DUIBA_COUPON":     {},
	"DUIBA_OBJECT":     {},
	"DUIBA_PHONEBILL":  {},
	"DUIBA_PHONEFLOW":  {},
	"DUIBA_VIRTUAL":    {},
	"DUIBA_GAME":       {},
	"DUIBA_HDTOOL":     {},
	"DUIBA_SIGN":       {},
	"DUIBA_REFUND":     {},
	"SYSTEM_REDUCE":    {},
	"SYSTEM_ADD":       {},
	"JHX":              {},
	"ARTICAL":          {},
	"RECOMMEND":        {},
	"DUIBA_POSTSALE":   {},
	"DUIBA_CANCELSHIP": {},
	"DUIBA_TASK":       {},
}
