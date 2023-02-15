package recyclemsg

import "encoding/json"

type IRecycleInfo interface {
	IUHDIOUHOIQWHIOEIOWEIOWEASOKA()
	JSON() ([]byte, error)
}
type RecycleInfo struct {
	Ch           string `json:"ch"`
	OrderNo      string `json:"orderNo"`
	MemberId     string `json:"memberId"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Number       string `json:"number"`
	CreateTime   string `json:"createTime"`
	CompleteTime string `json:"completeTime"`
	T            string `json:"t"`
	Sign         string `json:"sign"`
}

func (i RecycleInfo) IUHDIOUHOIQWHIOEIOWEIOWEASOKA() {}
func (i RecycleInfo) JSON() ([]byte, error) {
	return json.Marshal(i)
}

type RecycleFmyInfo struct {
	AppId          string         `json:"app_id"`
	NotificationAt string         `json:"notification_at"`
	Data           RecycleFmyData `json:"data"`
	Sign           string         `json:"sign"`
}
type RecycleFmyData struct {
	OrderSn          string `json:"order_sn"`
	Status           string `json:"status"`
	Weight           string `json:"weight"`
	Reason           string `json:"reason"`
	CourierRealName  string `json:"courier_real_name"`
	CourierPhone     string `json:"courier_phone"`
	CourierJobNumber string `json:"courier_job_number"`
	Waybill          string `json:"waybill"`
	Phone            string `json:"phone"`
}

func (i RecycleFmyInfo) IUHDIOUHOIQWHIOEIOWEIOWEASOKA() {}
func (i RecycleFmyInfo) JSON() ([]byte, error) {
	return json.Marshal(i)
}

type RecycleOolaInfo struct {
	Type                string `json:"type"`           //业务类型 1：回首订单成功
	OrderNo             string `json:"orderNo"`        //订单号，同类型同订单视为重复订单
	Name                string `json:"name"`           //type = 1，回收物品名称
	OolaUserId          string `json:"oolaUserId"`     //噢啦平台用户id
	ClientId            string `json:"clientId"`       //lvmiao用户id
	CreateTime          string `json:"createTime"`     //订单创建时间
	CompletionTime      string `json:"completionTime"` //订单完成时间
	BeanNum             string `json:"beanNum"`
	Sign                string `json:"sign"`                //加密串
	ProductCategoryName string `json:"productCategoryName"` //物品所属分类名称
	Qua                 string `json:"qua"`                 //用户下单时的数量&重量
	Unit                string `json:"unit"`                //与下单数量&重量关联的计量单位 如：公斤，个 等
}

func (i RecycleOolaInfo) IUHDIOUHOIQWHIOEIOWEIOWEASOKA() {}
func (i RecycleOolaInfo) JSON() ([]byte, error) {
	return json.Marshal(i)
}
