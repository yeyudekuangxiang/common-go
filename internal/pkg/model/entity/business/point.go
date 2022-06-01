package business

import "mio/internal/pkg/model"

type PointType string

const (
	PointTypeOnlineMeeting        PointType = "OnlineMeeting"        //线上会议
	PointTypeSaveWaterElectricity PointType = "SaveWaterElectricity" //节水节电
	PointTypePublicTransport      PointType = "PublicTransport"      //公交地铁
	PointTypeEvCar                PointType = "EvCar"                //电动车 电车充电
)

// Text 展示给用户看的
func (t PointType) Text() string {
	switch t {
	case PointTypeOnlineMeeting:
		return "线上会议"
	case PointTypeSaveWaterElectricity:
		return "节水节电"
	case PointTypePublicTransport:
		return "公交地铁"
	case PointTypeEvCar:
		return "电车充电"
	}
	return "未知类型"
}

// RealText 展示给管理员看的
func (t PointType) RealText() string {
	switch t {
	case PointTypeOnlineMeeting:
		return "线上会议"
	case PointTypeSaveWaterElectricity:
		return "节水节电"
	case PointTypePublicTransport:
		return "公交地铁"
	case PointTypeEvCar:
		return "电车充电"
	}
	return "未知类型"
}

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
