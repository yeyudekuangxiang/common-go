package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultCertificateService = CertificateService{repo: repository.DefaultCertificateRepository}

type CertificateService struct {
	repo repository.CertificateRepository
}

func (srv CertificateService) FindCertificate(by FindCertificateBy) (*entity.Certificate, error) {
	cert := srv.repo.FindCertificate(repository.FindCertificateBy{
		ProductItemId: by.ProductItemId,
		CertificateId: by.CertificateId,
	})
	return &cert, nil
}
