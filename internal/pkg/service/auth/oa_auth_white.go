package auth

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
)

var DefaultOaAuthWhiteService = OaAuthWhiteService{repo: repository.DefaultOaAuthWhiteRepository}

type OaAuthWhiteService struct {
	repo repository.OaAuthWhiteRepository
}

func (srv OaAuthWhiteService) FindBy(by FindOaAuthWhiteBy) (*entity.OaAuthWhite, error) {
	white := srv.repo.FindBy(repository.FindOaAuthWhiteBy{
		AppId:  by.AppId,
		Domain: by.Domain,
	})

	return &white, nil
}
