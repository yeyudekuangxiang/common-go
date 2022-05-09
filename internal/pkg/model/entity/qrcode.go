package entity

type QrCodeType string

const (
	QrCodeTypeSHARE       QrCodeType = "SHARE"
	QrCodeTypePARTNERSHIP QrCodeType = "PARTNERSHIP"
	QrCodeTypeGREENTORCH  QrCodeType = "GREEN_TORCH"
)

type QRCode struct {
	ID          int64      `json:"id"`
	QrCodeId    string     `json:"QrCodeId"`
	ImageUrl    string     `json:"imageUrl"`
	QrCodeType  QrCodeType `json:"qrCodeType"`
	OpenId      string     `json:"openId" gorm:"column:openid"`
	Description string     `json:"description"`
}

func (QRCode) TableName() string {
	return "qr_code"
}
