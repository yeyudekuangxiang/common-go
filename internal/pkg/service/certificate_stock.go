package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultCertificateStockService = CertificateStockService{repo: repository.DefaultCertificateStockRepository}

type CertificateStockService struct {
	repo repository.CertificateStockRepository
}

func (srv CertificateStockService) FindUnusedCertificate(certificateId string) (*entity.CertificateStock, error) {
	unusedStock, err := srv.repo.FindUnusedCertificate(certificateId)
	if err != nil {
		return nil, err
	}
	return &unusedStock, nil
}
