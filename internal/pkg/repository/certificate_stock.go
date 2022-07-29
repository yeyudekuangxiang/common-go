package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultCertificateStockRepository = CertificateStockRepository{DB: app.DB}

type CertificateStockRepository struct {
	DB *gorm.DB
}

func (repo CertificateStockRepository) FindUnusedCertificate(certificateId string) (entity.CertificateStock, error) {
	stock := entity.CertificateStock{}

	return stock, repo.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(stock).
			Set("gorm:query_option", "for update skip locked").
			Where("certificate_id = ? and used = false", certificateId).
			Take(&stock).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if stock.ID != 0 {
			stock.Used = true
			return tx.Save(&stock).Error
		}
		return nil
	})
}
func (repo CertificateStockRepository) Save(stock *entity.CertificateStock) error {
	return repo.DB.Save(stock).Error
}
