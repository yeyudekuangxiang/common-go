package model

type PointAddType string

const (
	PointAddTypeGame       PointAddType = "game"       //游戏
	PointAddTypeSign       PointAddType = "sign"       //签到
	PointAddTypeTask       PointAddType = "task"       //pk比赛
	PointAddTypeReSign     PointAddType = "reSign"     //补签
	PointAddTypePostSale   PointAddType = "postsale"   //售后退积分
	PointAddTypeCancelShip PointAddType = "cancelShip" //取消发货
	PointAddTypeHdTool     PointAddType = "hdtool"     //加积分活动
)

type PointAdd struct {
	Base
	Credits     IntStr       `json:"credits" form:"credits" alias:"credits"`
	Type        PointAddType `json:"type" form:"type" binding:"required" alias:"type"`
	OrderNum    string       `json:"orderNum" form:"orderNum" binding:"required" alias:"orderNum"`
	SubOrderNum string       `json:"subOrderNum" form:"subOrderNum" alias:"subOrderNum"`
	Description string       `json:"description" form:"description" alias:"description"`
	IP          string       `json:"ip" form:"ip" alias:"ip"`
}

func (p PointAdd) ToMap() map[string]string {
	m := p.Base.ToMap()
	m["credits"] = string(p.Credits)
	m["type"] = string(p.Type)
	m["orderNum"] = p.OrderNum
	m["subOrderNum"] = p.SubOrderNum
	m["description"] = p.Description
	m["ip"] = p.IP
	return m
}
