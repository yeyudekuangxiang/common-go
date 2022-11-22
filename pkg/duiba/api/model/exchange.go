package model

import (
	"strconv"
)

type ExchangeType string

const (
	ExchangeTypeAlipay    ExchangeType = "alipay"    //支付宝
	ExchangeTypeQB        ExchangeType = "qb"        //Q币
	ExchangeTypeCoupon    ExchangeType = "coupon"    //优惠券
	ExchangeTypeObject    ExchangeType = "object"    //实物
	ExchangeTypePhoneBill ExchangeType = "phonebill" //话费
	ExchangeTypePhoneFlow ExchangeType = "phoneflow" //流量
	ExchangeTypeVirtual   ExchangeType = "virtual"   //虚拟商品
	ExchangeTypeGame      ExchangeType = "game"      //游戏
	ExchangeTypeHdTool    ExchangeType = "hdtool"    //活动抽奖
	ExchangeTypeHdSign    ExchangeType = "sign"      //签到
)

// Exchange 扣积分接口 https://www.duiba.com.cn/tech_doc_book/server/servers/consume_credits_api_2.html#api%E6%96%87%E6%A1%A3
type Exchange struct {
	Base
	Credits     int64        `json:"credits" form:"credits"`                                                //本次兑换扣除的积分
	ItemCode    string       `json:"itemCode" form:"itemCode"`                                              //自有商品商品编码(非必须字段)
	Description string       `json:"description" form:"description" binding:"required" alias:"description"` //本次积分消耗的描述(带中文，请用utf-8进行url解码)
	OrderNum    string       `json:"orderNum" form:"orderNum" binding:"required" alias:"orderNum"`          //兑吧订单号(请记录到数据库中)
	Type        ExchangeType `json:"type" form:"type" binding:"required" alias:"type"`                      //兑换类型
	FacePrice   int          `json:"facePrice" form:"facePrice"`                                            //兑换商品的市场价值，单位是分，请自行转换单位
	ActualPrice int          `json:"actualPrice" form:"actualPrice"`                                        //此次兑换实际扣除开发者账户费用，单位为分
	IP          string       `json:"ip" form:"ip"`                                                          //用户ip，不保证获取到
	QQ          string       `json:"qq" form:"qq"`                                                          //直冲商品Q币商品，QQ号码回传参数，其他商品不传该参数
	Phone       string       `json:"phone" form:"phone"`                                                    //直冲类话费商品手机号回传参数，非话费商品不传该参数
	Alipay      string       `json:"alipay" form:"alipay"`                                                  //支付宝充值商品支付宝账号参数回传，非支付宝商品不传该参数
	WaitAudit   bool         `json:"waitAudit" form:"waitAudit"`                                            //是否需要审核(如需在自身系统进行审核处理，请记录下此信息)
	Params      string       `json:"params" form:"params"`                                                  //详情参数，不同的类型，请求时传不同的内容，中间用英文冒号分隔。(支付宝类型带中文，请用utf-8进行解码) 实物商品：返回收货信息(姓名:手机号:省份:城市:区域:街道:详细地址)、支付宝：返回账号信息(支付宝账号:实名)、话费：返回手机号、QB：返回QQ号
}

func (form Exchange) ToMap() map[string]string {
	return map[string]string{
		"uid":         form.Uid,
		"credits":     strconv.FormatInt(form.Credits, 10),
		"itemCode":    form.ItemCode,
		"appKey":      form.AppKey,
		"timestamp":   string(form.Timestamp),
		"description": form.Description,
		"orderNum":    form.OrderNum,
		"type":        string(form.Type),
		"facePrice":   strconv.Itoa(form.FacePrice),
		"actualPrice": strconv.Itoa(form.ActualPrice),
		"ip":          form.IP,
		"qq":          form.QQ,
		"phone":       form.Phone,
		"alipay":      form.Alipay,
		"waitAudit":   strconv.FormatBool(form.WaitAudit),
		"params":      form.Params,
		"sign":        form.Sign,
	}
}

// ExchangeResult 兑换结果通知 https://www.duiba.com.cn/tech_doc_book/server/servers/notify_api_2.html
type ExchangeResult struct {
	Base
	Success      bool   `json:"success" form:"success"`                                       //兑换是否成功，状态是true和false
	ErrorMessage string `json:"errorMessage" form:"errorMessage"`                             //出错原因(带中文，请用utf-8进行解码)
	OrderNum     string `json:"orderNum" form:"orderNum" binding:"required" alias:"orderNum"` //兑吧订单号(请记录到数据库中)
	BizId        string `json:"bizId" form:"bizId"`                                           //开发者的订单号
}

func (form ExchangeResult) ToMap() map[string]string {
	return map[string]string{
		"uid":          form.Uid,
		"appKey":       form.AppKey,
		"timestamp":    string(form.Timestamp),
		"success":      strconv.FormatBool(form.Success),
		"errorMessage": form.ErrorMessage,
		"orderNum":     form.OrderNum,
		"bizId":        form.BizId,
		"sign":         form.Sign,
	}
}
