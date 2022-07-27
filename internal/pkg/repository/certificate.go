package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCertificateRepository = CertificateRepository{DB: app.DB}

type CertificateRepository struct {
	DB *gorm.DB
}

func (repo CertificateRepository) FindCertificate(by FindCertificateBy) entity.Certificate {
	cert := entity.Certificate{}
	db := repo.DB.Model(cert)

	if by.ProductItemId != "" {
		db.Where("product_item_id = ?", by.ProductItemId)
	}
	if by.CertificateId != "" {
		db.Where("certificate_id = ?", by.CertificateId)
	}

	err := db.Take(&cert).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	return cert
}
