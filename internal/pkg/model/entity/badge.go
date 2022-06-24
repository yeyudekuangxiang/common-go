package entity

import "mio/internal/pkg/model"

type Badge struct {
	ID            int64           `json:"id" gorm:"primaryKey;type:serial8;not null;comment:用户证书记录"`
	Code          string          `json:"code" gorm:"type:varchar(100);not null;default:'';comment:证书编号"`
	OpenId        string          `json:"openId" gorm:"column:openid;type:varchar(100);not null;comment:用户openid"`
	CertificateId string          `json:"certificateId" gorm:"type:varchar(100);not null;comment:证书id"`
	ProductItemId string          `json:"productItemId" gorm:"type:varchar(100);not null;comment:对应商品id"`
	CreateTime    model.Time      `json:"createTime" gorm:"type:timestamptz;not null;comment:创建时间"`
	Partnership   PartnershipType `json:"partnership" gorm:"type:varchar(100);not null;default:'';comment:合作伙伴"`
	OrderId       string          `json:"orderId" gorm:"type:varchar(100);not null;default:'';comment:订单编号"`
	ImageUrl      string          `json:"imageUrl" gorm:"type:varchar(1000);not null;default:'';comment:证书图片"`
	IsNew         bool            `json:"isNew" gorm:"type:bool;not null;default:false;comment:是否是新获得"`
}

func (Badge) TableName() string {
	return "badge"
}
