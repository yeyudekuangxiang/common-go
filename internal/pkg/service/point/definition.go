package point

import "mio/internal/pkg/model/entity"

type commandDescription struct {
	Times      int64                            //限制次数 0表示不限次数
	Amount     int64                            //积分
	MaxAmount  int64                            //每日最多获取积分
	Fn         func(*defaultClientHandle) error //执行方法
	FnPageData func(*defaultClientHandle) (map[string]interface{}, error)
	IsOpen     bool //是否实现
}

type CollectType string
type CollectRules map[entity.PointTransactionType][]string
type RuleTranslate map[entity.PointTransactionType]map[string]string

var rules = CollectRules{
	"COFFEE_CUP":    []string{"自带杯", "单号", "订单"},
	"BIKE_RIDE":     []string{"骑行", "单车", "骑车", "bike", "出行", "哈啰", "摩拜", "青桔"},
	"POWER_REPLACE": []string{"订单编号", "已支付"},
}

var identifyChRules = CollectRules{
	"POWER_REPLACE": []string{"订单编号", "度", "支付状态", "充电量", "已支付"},
}

var identifyEnRules = RuleTranslate{
	"POWER_REPLACE": map[string]string{
		"订单编号": "orderId",
		"度":    "kwh",
	},
}

//方法定义
var commandMap = map[string]*commandDescription{
	//times:每天可以提交次数 amount:1度电10积分 maxAmount:每天最多获取积分
	"POWER_REPLACE": {Fn: (*defaultClientHandle).powerReplace, Times: 1, Amount: 10, MaxAmount: 300},
}

//before弹框页面获取数据
var pageDataMap = map[string]*commandDescription{
	"POWER_REPLACE":    {FnPageData: (*defaultClientHandle).powerReplacePageData},    //换电
	"OOLA_RECYCLE":     {FnPageData: (*defaultClientHandle).oolaRecyclePageData},     //oola旧物回收
	"FAST_ELECTRICITY": {FnPageData: (*defaultClientHandle).fastElectricityPageData}, //快电
	"FMY_RECYCLE":      {FnPageData: (*defaultClientHandle).fmyRecyclePageData},
	"REDUCE_PLASTIC":   {FnPageData: (*defaultClientHandle).reducePlasticPageData},
	"JHX":              {FnPageData: (*defaultClientHandle).jhxPageData},
}
