package entity

import "mio/internal/pkg/model"

type DuiBaActivity struct {
	ID                  int
	ActivityId          string
	ActivityUrl         string
	AfterEndActivity    int //1 跳转到兑吧商城首页 2弹出提示活动未开始并且返回首页
	BeforeStartActivity int //1 跳转到兑吧商城首页 2弹出提示活动已结束并且返回首页
	StartTime           model.Time
	EndTime             model.Time
	CreatedAt           model.Time
	UpdatedAt           model.Time
	DeletedAt           model.Time
	RiskLimit           int //用户风险等级限制
}

func (DuiBaActivity) TableName() string {
	return "duiba_activity"
}
