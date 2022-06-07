package business

import (
	"mio/internal/pkg/model"
)

type PointCollectHistory struct {
	ID            int64      `json:"id" gorm:"primaryKey;not null;type:serial8;comment:积分获取信息记录表"`
	TransactionId string     `json:"transactionId" gorm:"not null;type:varchar(100);comment:过程id"`
	BUserId       int64      `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	Type          PointType  `json:"type" gorm:"not null;type:varchar(50);comment:积分类型"`
	Info          string     `json:"info" gorm:"not null;type:varchar(1000);default:'';comment:附带信息 json object 同一个type的info格式必须统一"`
	TimePoint     model.Time `json:"timePoint" gorm:"not null;type:timestamptz;comment:时间点 2006-01-02 00:00:00"`
	CreatedAt     model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt     model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (PointCollectHistory) TableName() string {
	return "business_point_collect_history"
}
