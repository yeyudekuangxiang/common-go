package business

import "mio/internal/pkg/model"

type PointTransaction struct {
	ID        int64      `json:"id" gorm:"primaryKey;not null;type:serial8;comment:积分变动表"`
	BUserId   int64      `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	Value     int        `json:"value" gorm:"not null;type:int4;comment:积分变动数量"`
	Type      string     `json:"type" gorm:"not null;type:varchar(20);comment:积分类型"`
	OrderId   string     `json:"orderId" gorm:"not null;type:varchar(255);default:'';comment:相关订单id"`
	Info      string     `json:"info" gorm:"not null;type:varchar(1000);default:'';comment:附带信息 json object 同一个type的info格式必须统一"`
	CreatedAt model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (PointTransaction) TableName() string {
	return "business_point_transaction"
}
