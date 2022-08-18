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
	IsPhone             DuiBaActivityIsPhone
	VipType             DuiBaActivityVipType
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

type DuiBaActivityIsPhone int8

const (
	DuiBaActivityIsPhoneYes DuiBaActivityIsPhone = 1
	DuiBaActivityIsPhoneNo  DuiBaActivityIsPhone = 2
)

type DuiBaActivityVipType int8

const (
	DuiBaActivityVipTypeNewUser             DuiBaActivityVipType = 1
	DuiBaActivityIsPhoneAnniversaryActivity DuiBaActivityVipType = 2
)

const (
	DuiBaActivityRiskLimitZero  int = 0
	DuiBaActivityRiskLimitOne   int = 1
	DuiBaActivityRiskLimitTwo   int = 2
	DuiBaActivityRiskLimitThree int = 3
	DuiBaActivityRiskLimitFour  int = 4
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
	DuiBaActivityIsPhoneMap = map[DuiBaActivityIsPhone]string{
		DuiBaActivityIsPhoneYes: "需要绑定",
		DuiBaActivityIsPhoneNo:  "不需要绑定",
	}
	DuiBaActivityVipTypeMap = map[DuiBaActivityVipType]string{
		DuiBaActivityVipTypeNewUser:             "新用户",
		DuiBaActivityIsPhoneAnniversaryActivity: "周年庆活动",
	}
	DuiBaActivityDuiBaActivityMap = map[int]string{
		DuiBaActivityRiskLimitZero:  "风险等级0",
		DuiBaActivityRiskLimitOne:   "风险等级1",
		DuiBaActivityRiskLimitTwo:   "风险等级2",
		DuiBaActivityRiskLimitThree: "风险等级3",
		DuiBaActivityRiskLimitFour:  "风险等级4",
	}
)

func (DuiBaActivity) TableName() string {
	return "duiba_activity"
}
