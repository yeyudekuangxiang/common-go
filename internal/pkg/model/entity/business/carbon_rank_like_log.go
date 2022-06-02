package business

import "mio/internal/pkg/model"

type CarbonRankLikeLog struct {
	ID         int64          `json:"id" gorm:"primaryKey;not null;type:serial8;comment:碳减排排行榜点赞记录"`
	DateType   RankDateType   `json:"dateType" gorm:"not null;type:varchar(20);comment:点赞榜单类型 day(日榜) week(周榜) month(月榜)"`
	ObjectType RankObjectType `json:"objectType" gorm:"not null;type:varchar(20);comment:点赞对象类型 user(用户) department(部门)"`
	BUserId    int64          `json:"-" gorm:"not null;type:int8;comment:点赞人id 企业用户表主键"`
	Pid        int64          `json:"pid" gorm:"not null;type:int8;comment:根据object_type 分别对应 user表和department表主键"`
	Status     int8           `json:"status" gorm:"not null;type:int2;comment:点赞状态 1 已点赞 2取消点赞"`
	TimePoint  model.Time     `json:"timePoint" gorm:"not null;type:varchar(20);comment:根据date_type 记录日 周 月 开始的时间节点 2006-01-02 00:00:00"` //同一个用户  日榜每天都会有一条记录 周榜每周有一条记录 月榜每月有一条记录
	CreatedAt  model.Time     `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt  model.Time     `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (CarbonRankLikeLog) TableName() string {
	return "business_carbon_rank_like_log"
}
