package business

import (
	"github.com/shopspring/decimal"
	"mio/internal/pkg/model"
	"mio/internal/pkg/util/timeutils"
	"time"
)

type RankDateType string

const (
	RankDateTypeDay   RankDateType = "day"
	RankDateTypeWeek  RankDateType = "week"
	RankDateTypeMonth RankDateType = "month"
)

// ParseLastTime 昨天开始和结束 上周开始和结束 上个月开始和结束
func (rdt RankDateType) ParseLastTime() (time.Time, time.Time) {
	t := timeutils.Now()
	var start, end time.Time
	switch rdt {
	case RankDateTypeDay:
		start = t.AddDay(-1).StartOfDay().Time
		end = t.AddDay(-1).EndOfDay().Time
	case RankDateTypeWeek:
		start = t.AddWeek(-1).StartOfWeek().Time
		end = t.AddWeek(-1).EndOfWeek().Time
	case RankDateTypeMonth:
		start = t.AddMonth(-1).StartOfMonth().Time
		end = t.AddMonth(-1).EndOfMonth().Time
	}
	return start, end
}

type RankObjectType string

const (
	RankObjectTypeUser       RankObjectType = "user"
	RankObjectTypeDepartment RankObjectType = "department"
)

type CarbonRank struct {
	ID         int64           `json:"id" gorm:"primaryKey;not null;type:serial8;comment:碳排行榜"`
	DateType   RankDateType    `json:"dateType" gorm:"not null;type:varchar(20);comment:点赞榜单类型 day(日榜) week(周榜) month(月榜)"`
	ObjectType RankObjectType  `json:"objectType" gorm:"not null;type:varchar(20);comment:点赞对象类型 user(用户) department(部门)"`
	Value      decimal.Decimal `json:"value" gorm:"not null;type:decimal(20,2);comment:周期内获取到的碳积分数量"`
	Rank       int             `json:"rank" gorm:"not null;type:int4;comment:榜单排名"`
	Pid        int64           `json:"pid" gorm:"not null;type:int8;comment:根据object_type 分别对应 user表和department表主键"`
	LikeNum    int             `json:"likeNum" gorm:"not null;type:int4;comment:点赞数量"`
	TimePoint  model.Time      `json:"timePoint" gorm:"not null;type:timestamp;comment:根据date_type 记录日 周 月 开始的时间节点 2006-01-02 00:00:00"` //同一个用户  日榜每天都会有一条记录 周榜每周有一条记录 月榜每月有一条记录
	CreatedAt  model.Time      `json:"createdAt" gorm:"not null;type:timestamp"`
	UpdatedAt  model.Time      `json:"updatedAt" gorm:"not null;type:timestamp"`
}

func (CarbonRank) TableName() string {
	return "business_carbon_rank"
}
