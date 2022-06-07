package business

import (
	ebusiness "mio/internal/pkg/model/entity/business"
	rbusiness "mio/internal/pkg/repository/business"
)

var DefaultCarbonSceneService = CarbonSceneService{repo: rbusiness.DefaultCarbonSceneRepository}

type CarbonSceneService struct {
	repo rbusiness.CarbonSceneRepository
}

func (srv CarbonSceneService) FindScene(t ebusiness.CarbonType) (*ebusiness.CarbonScene, error) {
	scene := srv.repo.FindScene(t)
	return &scene, nil
}
