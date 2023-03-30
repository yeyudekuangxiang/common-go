package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultBannerController = BannerController{}

type BannerController struct {
}

func (BannerController) GetBannerList(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetGetBannerListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	var displays []string
	if form.Display == "" {
		displays = append(displays, "min", "all")
	} else {
		displays = append(displays, form.Display, "all")
	}
	list, err := service.DefaultBannerService.List(srv_types.GetBannerListDTO{
		Scene:    entity.BannerScene(form.Scene),
		Status:   1,
		Displays: displays,
	})
	if err != nil {
		return nil, err
	}

	bannerVOList := make([]api_types.BannerVO, 0)
	for _, banner := range list {
		bannerVOList = append(bannerVOList, api_types.BannerVO{
			ID:       banner.ID,
			Name:     banner.Name,
			ImageUrl: banner.ImageUrl,
			Type:     string(banner.Type),
			Redirect: banner.Redirect,
			AppId:    banner.AppId,
			Ext:      banner.Ext,
		})
	}

	return gin.H{"list": bannerVOList}, nil
}
