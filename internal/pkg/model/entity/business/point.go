package business

import "mio/internal/pkg/model"

type PointType string

const (
	PointTypeOnlineMeeting        CarbonType = "OnlineMeeting"        //线上会议
	PointTypeSaveWaterElectricity CarbonType = "SaveWaterElectricity" //节水节电
	PointTypePublicTransport      CarbonType = "PublicTransport"      //公交地铁
	PointTypeEvCar                CarbonType = "EvCar"                //电动车 电车充电
)

type Point struct {
	ID        int64      `json:"id" gorm:"primaryKey;not null;type:serial8;comment:积分账户表"`
	BUserId   int64      `json:"-" gorm:"not null;type:int8;comment:企业用户表主键"`
	Point     int64      `json:"point" gorm:"not null;type:int8;comment:积分余额"`
	UsedPoint int64      `json:"usedPoint" gorm:"not null;type:int8;comment:已使用的积分数量"`
	CreatedAt model.Time `json:"createdAt" gorm:"not null;type:timestamptz"`
	UpdatedAt model.Time `json:"updatedAt" gorm:"not null;type:timestamptz"`
}

func (Point) TableName() string {
	return "business_point"
}
