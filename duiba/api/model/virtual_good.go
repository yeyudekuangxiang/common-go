package model

// VirtualGood 虚拟商品充值 https://docs.duiba.com.cn/tech_doc_book/optional/virtual_2.html
type VirtualGood struct {
	Base
	OrderNum     string `json:"orderNum" form:"orderNum" binding:"required" alias:"兑吧订单号"`
	DevelopBizId string `json:"developBizId" form:"developBizId" `
	Params       string `json:"params" form:"params" binding:"required" alias:"虚拟商品标识符"`
	Description  string `json:"description" form:"description" `
	Account      string `json:"account" form:"account" `
}

func (v VirtualGood) ToMap() map[string]string {
	m := v.Base.ToMap()
	m["developBizId"] = v.DevelopBizId
	m["orderNum"] = v.OrderNum
	m["params"] = v.Params
	m["description"] = v.Description
	m["account"] = v.Account
	return m
}
