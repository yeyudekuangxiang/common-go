package service

import (
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
)

var DefaultBannerService = BannerService{repo: repository.DefaultBannerRepository}

type BannerService struct {
	repo repository.BannerRepository
}

func (srv BannerService) List(dto srv_types.GetBannerListDTO) ([]entity.Banner, error) {
	bannerDo := repotypes.GetBannerListDO{OrderBy: entity.OrderByList{entity.OrderByBannerSortAsc}}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return nil, err
	}
	return srv.repo.List(bannerDo)
}
