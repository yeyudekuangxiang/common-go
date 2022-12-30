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

var recyclePointOfName = map[int]interface{}{
	1: 21.0, //1000g : 21 积分
	2: 6.0,
	3: map[string]float64{
		"默认":  113.0,
		"手机":  113.0,
		"平板":  409.0,
		"一体机": 1031.0,
		"笔记本": 1911.0,
	},

	4: map[string]float64{
		"默认":  69.0,
		"电视":  205,
		"冰箱":  384.0,
		"空调":  69.0,
		"洗衣机": 690.0,
	},
	5: map[string]float64{
		"默认":   13.8,
		"柜子":   13.8,
		"桌椅":   104.88,
		"健身器材": 376.28,
		"床":    692.3,
		"沙发":   230,
		"茶几":   1863,
		"床垫":   532.8,
		"综合式":  13.8, //柜子
		"架":    13.8,
		"台":    104.88, //桌椅
		"机":    376.28,
		"车":    376.28,
	},
	6:   1000.0,
	7:   10.0,
	100: 100.0,
}

var recycleCo2OfName = map[int]interface{}{
	1: 4500.0,
	2: 140.0,
	3: map[string]float64{
		"默认":  25000,
		"手机":  25000,
		"平板":  89000,
		"一体机": 224000,
		"笔记本": 415000,
	},
	4: map[string]float64{
		"默认":  15000,
		"电视":  45000,
		"冰箱":  83000,
		"空调":  15000,
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
		"综合式":  300, //柜子
		"架":    300,
		"台":    22800, //桌椅
		"机":    8180,
		"车":    8180,
	},
	6:   0.0,
	7:   2420.0,
	100: 8966.8,
}

var recycleMaxPoint = map[int]int{
	1:   2709, //1000g : 21 积分
	2:   322,
	3:   1911,
	4:   920,
	5:   1863,
	6:   1000,
	7:   969,
	100: 7500,
}
