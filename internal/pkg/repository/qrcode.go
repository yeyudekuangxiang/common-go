package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

type QrCodeRepository struct {
	db *gorm.DB
}

func NewQrCodeRepository(db *gorm.DB) *QrCodeRepository {
	return &QrCodeRepository{db: db}
}

//GetQrCode 获取二维码
// scene entity.QrCodeScene 二维码的使用场景
// key  二维码的key值 key值和scene应该组成唯一索引
// 返回值
// qr *entity.QRCode 二维码信息
// exist bool 二维码信息是否存在 true代表存在 false代表不存在
// err error 查询异常错误信息
func (repo QrCodeRepository) GetQrCode(key string) (qr *entity.QRCode, exist bool, err error) {
	qr = &entity.QRCode{}
	err = app.DB.Where("key = ?", key).First(qr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return qr, true, nil
}

func (repo QrCodeRepository) Create(code *entity.QRCode) error {
	return repo.db.Create(code).Error
}
