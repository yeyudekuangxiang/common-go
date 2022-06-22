package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/service_types"
)

var DefaultUploadSceneService = UploadSceneService{repo: repository.DefaultUploadSceneRepository}

type UploadSceneService struct {
	repo repository.UploadSceneRepository
}

func (srv UploadSceneService) FindUploadScene(param service_types.FindSceneParam) (*entity.UploadScene, error) {
	return srv.repo.FindScene(repository.FindSceneBy{
		Scene: param.Scene,
	})
}
