package business

import (
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCompanyCarbonSceneService = CompanyCarbonSceneService{repo: rbusiness.DefaultCompanyCarbonSceneRepository}

type CompanyCarbonSceneService struct {
	repo rbusiness.CompanyCarbonSceneRepository
}

func (srv CompanyCarbonSceneService) FindCompanyScene(param FindCompanyCarbonSceneParam) (*ebusiness.CompanyCarbonScene, error) {
	scene := srv.repo.FindCompanyScene(rbusiness.FindCompanyCarbonSceneBy{
		CompanyId:     param.CompanyId,
		CarbonSceneId: param.CarbonSceneId,
	})
	return &scene, nil
}
