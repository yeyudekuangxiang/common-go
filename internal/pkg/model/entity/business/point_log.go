package business

import (
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
)

const (
	OrderByPointLogCTDESC entity.OrderBy = "order_by_business_point_ct_desc" //创建时间倒叙
)

type PointLog struct {
	ID            int64         `json:"id" gorm:"primaryKey;not null;type:serial8;comment:积分变动表"`
	TransactionId string        `json:"transactionId" gorm:"not null;type:varchar(100);comment:过程id"`
	BUserId       int64         `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	Value         int           `json:"value" gorm:"not null;type:int4;comment:积分变动数量"`
	Type          PointType     `json:"type" gorm:"not null;type:varchar(50);comment:积分类型"`
	OrderId       string        `json:"orderId" gorm:"not null;type:varchar(255);default:'';comment:相关订单id"`
	Info          PointTypeInfo `json:"info" gorm:"not null;type:varchar(1000);default:'';comment:附带信息 json object 同一个type的info格式必须统一"`
	CreatedAt     model.Time    `json:"createdAt" gorm:"not null;type:timestamp"`
	UpdatedAt     model.Time    `json:"updatedAt" gorm:"not null;type:timestamp"`
}

func (PointLog) TableName() string {
	return "business_point_log"
}
