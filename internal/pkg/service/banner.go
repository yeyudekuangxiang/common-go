package service

import (
	"github.com/pkg/errors"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"strings"
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

func (srv BannerService) Create(dto srv_types.CreateBannerDTO) error {
	//判断名称和图片是否存在
	banner, err := srv.repo.GetExistOne(repotypes.GetBannerExistDO{Name: dto.Name, ImageUrl: dto.ImageUrl})
	if err != nil {
		return err
	}
	if banner.ID != 0 {
		return errors.New("banner名称或图片已存在")
	}
	if dto.Type == entity.BannerTypeH5 {
		dto.Redirect = "pages/load_bearing/webview/index?url=" + dto.Redirect
	}
	bannerDo := entity.Banner{
		CreateTime: model.NewTime(),
		UpdateTime: model.NewTime()}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return err
	}
	return srv.repo.Create(&bannerDo)
}

func (srv BannerService) Update(dto srv_types.UpdateBannerDTO) error {
	bannerOne := srv.repo.GetById(dto.Id)
	if bannerOne.ID == 0 {
		return errors.New("banner记录不存在")
	}
	//判断名称和图片是否存在
	banner, err := srv.repo.GetExistOne(repotypes.GetBannerExistDO{Name: dto.Name, NotId: dto.Id})
	if err != nil {
		return err
	}
	if banner.ID != 0 {
		return errors.New("banner名称已存在")
	}

	if dto.Type == entity.BannerTypeH5 {
		if find := strings.Contains(dto.Redirect, "pages/load_bearing/webview/index?url="); !find {
			dto.Redirect = "pages/load_bearing/webview/index?url=:" + dto.Redirect
		}
	}
	bannerDo := entity.Banner{
		UpdateTime: model.NewTime()}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return err
	}
	return srv.repo.Save(&bannerDo)
}

func (srv BannerService) GetBannerPageList(dto srv_types.GetPageBannerDTO) ([]entity.Banner, int64, error) {
	bannerDo := repotypes.GetBannerPageDO{OrderBy: entity.OrderByList{entity.OrderByBannerSortAsc}}
	if err := util.MapTo(dto, &bannerDo); err != nil {
		return nil, 0, err
	}
	list, total, err := srv.repo.Page(bannerDo)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
