package entity

import "time"

// QrCodeScene 使用场景
type QrCodeScene string

const (
	QrCodeSceneInvite     QrCodeScene = "invite"      //邀请得积分
	QrCodeSceneTopicShare QrCodeScene = "topic_share" //酷喵圈分享
)

// QrCodeSource 二维码来源
type QrCodeSource string

const (
	QrCodeSourceWeappUnlimited QrCodeSource = "weapp_unlimited" //微信小程序无数量限制小程序码
	QrCodeSourceWeappLimited   QrCodeSource = "weapp_limited"   //微信小程序有数量限制小程序吗
	QrCodeSourceWeappQr        QrCodeSource = "weapp_qrcode"    //微信小程序有数量限制二维码
	QrCodeSourceCommon         QrCodeSource = "common"          //普通二维码
)

type QRCode struct {
	ID           int64        `json:"id"`
	QrCodeId     string       `json:"QrCodeId"`
	ImagePath    string       `json:"imagePath"`
	QrCodeScene  QrCodeScene  `json:"qrCodeScene"`                 //使用场景
	QrCodeSource QrCodeSource `json:"qrCodeSource"`                //来源
	Key          string       `json:"key"`                         //key 和 QrCodeScene 组成唯一索引
	Content      string       `json:"content"`                     //内容
	Ext          string       `json:"ext"`                         //额外的参数
	OpenId       string       `json:"openId" gorm:"column:openid"` //用于记录生成用户的id
	Description  string       `json:"description"`
	CreatedAt    time.Time    `json:"createdAt"` //创建时间
}

func (QRCode) TableName() string {
	return "qr_code"
}
