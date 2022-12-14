package recycle

type FmySignParams struct {
	AppId          string         `json:"app_id" form:"app_id" binding:"required"`
	NotificationAt string         `json:"notification_at" form:"notification_at" binding:"required"`
	Data           recycleFmyData `json:"data" form:"data"`
	Sign           string         `json:"sign,omitempty"`
}

type recycleFmyData struct {
	OrderSn          string `json:"order_sn,omitempty" binding:"required"`
	Status           string `json:"status,omitempty" binding:"required"`
	Weight           string `json:"weight,omitempty"`
	Reason           string `json:"reason,omitempty"`
	CourierRealName  string `json:"courier_real_name,omitempty"`
	CourierPhone     string `json:"courier_phone,omitempty"`
	CourierJobNumber string `json:"courier_job_number,omitempty"`
	Waybill          string `json:"waybill,omitempty"`
	Phone            string `json:"phone,omitempty"`
}

type oolaSignParams struct {
	Type                int    `json:"type" form:"type" binding:"required"`                     //业务类型 1：回首订单成功
	OrderNo             string `json:"orderNo" form:"orderNo" binding:"required"`               //订单号，同类型同订单视为重复订单
	Name                string `json:"name" form:"name" binding:"required"`                     //type = 1，回收物品名称
	ProductCategoryName string `json:"productCategoryName" form:"productCategoryName"`          //物品所属分类名称
	Qua                 string `json:"qua" form:"qua"`                                          //用户下单时的数量&重量
	Unit                string `json:"unit" form:"unit"`                                        //与下单数量&重量关联的计量单位 如：公斤，个 等
	OolaUserId          int    `json:"oolaUserId" form:"oolaUserId" binding:"required"`         //噢啦平台用户id
	ClientId            string `json:"clientId" form:"clientId" binding:"required"`             //lvmiao用户id
	CreateTime          string `json:"createTime" form:"createTime" binding:"required"`         //订单创建时间
	CompletionTime      string `json:"completionTime" form:"completionTime" binding:"required"` //订单完成时间
	//Sign                string `json:"sign" form:"sign" binding:"required"`                     //加密串
}

// 回收 台 单位对应积分 比如 电视机 1台 获得 69积分
var recyclePointV2 = map[int]float64{
	1: 21, //1000g : 21 积分
	2: 6,  //1000g : 6 积分
	//3:   113,
	//4:   384,
	6:   1000,
	100: 100,
}

// 回收 台/重量 单位对应减碳量 比如 电视机 1台 获得 15000g 减碳量
var recycleCo2V2 = map[int]float64{
	1: 4500, //1000g : 4500g
	2: 1400, //1000g : 1400g
	//3:   25000,
	//4:   83000,
	6:   0,
	100: 8966.8,
}

var recyclePointForName = map[int]interface{}{
	1: 21, //1000g : 21 积分
	2: 6,
	3: map[string]float64{
		"默认": 113,

		"手机":  113,
		"平板":  409,
		"一体机": 1031,
		"笔记本": 1911,
	},

	4: map[string]float64{
		"默认":  69,
		"电视":  69,
		"冰箱":  384,
		"空调":  205,
		"洗衣机": 690,
	},
	5: map[string]float64{
		"默认":   13.8,
		"柜子":   13.8,
		"桌椅":   104.88,
		"健身器材": 376.28,
		"床":    692.3,
		"沙发":   1863,
		"茶几":   230,
		"床垫":   532.8,
	},
	6: 1000,
	7: map[string]float64{
		"默认":    5.3,
		"塑料":    14.35,
		"金属":    14.26,
		"纸类":    6,
		"衣物":    5.3,
		"软塑包装":  14.9,
		"复合纸包装": 4.9,
	},
	100: 100,
}

var recycleCo2ForName = map[int]interface{}{
	1: 4500, //1000g : 4500g
	2: 1400,
	3: map[string]float64{
		"默认":  25000,
		"手机":  25000,
		"平板":  89000,
		"一体机": 224000,
		"笔记本": 415000,
	},
	4: map[string]float64{
		"默认":  15000,
		"电视":  15000,
		"冰箱":  83000,
		"空调":  45000,
		"洗衣机": 150000,
	},
	5: map[string]float64{
		"默认":   300,
		"柜子":   300,
		"桌椅":   22800,
		"健身器材": 8180,
		"床":    301000,
		"沙发":   100000,
		"茶几":   810000,
		"床垫":   236000,
	},
	6: 0,
	7: map[string]float64{
		"默认":    140,
		"塑料":    3120,
		"金属":    3100,
		"纸类":    140,
		"衣物":    4000,
		"软塑包装":  3240,
		"复合纸包装": 890,
	},
	100: 8966.8,
}
