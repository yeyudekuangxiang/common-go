package entity

import (
	"time"
)

type DuiBaActivity struct {
	ID                  int
	ActivityId          string
	ActivityUrl         string
	AfterEndActivity    int //1 跳转到兑吧商城首页 2弹出提示活动未开始并且返回首页
	BeforeStartActiviry int //1 跳转到兑吧商城首页 2弹出提示活动已结束并且返回首页
	StartTime           time.Time
	EndTime             time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time
	RiskLimit           int //用户风险等级限制
	Cid                 int64
	Type                DuiBaActivityType
	IsShare             DuiBaActivityIsShare
	Name                string
	Status              DuiBaActivityStatus
}

type DuiBaActivityType int8

const (
	DuiBaActivityGoodsShow DuiBaActivityType = 1
	DuiBaActivityLandPage  DuiBaActivityType = 2
)

type DuiBaActivityIsShare int8

const (
	DuiBaActivityIsShareYes DuiBaActivityIsShare = 1
	DuiBaActivityIsShareNo  DuiBaActivityIsShare = 2
)

type DuiBaActivityStatus int8

const (
	DuiBaActivityStatusYes DuiBaActivityStatus = 1
	DuiBaActivityStatusNo  DuiBaActivityStatus = 2
)

var (
	DuiBaActivityTypeMap = map[DuiBaActivityType]string{
		DuiBaActivityGoodsShow: "商品详情",
		DuiBaActivityLandPage:  "落地页链接",
	}
	DuiBaActivityIsShareMap = map[DuiBaActivityIsShare]string{
		DuiBaActivityIsShareYes: "是",
		DuiBaActivityIsShareNo:  "不是",
	}
	DuiBaActivityStatusMap = map[DuiBaActivityStatus]string{
		DuiBaActivityStatusYes: "正常",
		DuiBaActivityStatusNo:  "删除",
	}
)

func (DuiBaActivity) TableName() string {
	return "duiba_activity"
}
