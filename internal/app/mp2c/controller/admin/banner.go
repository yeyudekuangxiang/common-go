package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultBannerController = NewBannerController(service.DefaultBannerService)

func NewBannerController(channel service.BannerService) BannerController {
	return BannerController{service: channel}
}

type BannerController struct {
	service service.BannerService
}

/*创建banner*/
func (ctl BannerController) Create(c *gin.Context) (gin.H, error) {
	var form CreateBannerForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultBannerService.Create(srv_types.CreateBannerDTO{
		Name:     form.Name,
		ImageUrl: form.ImageUrl,
		Scene:    form.Scene,
		Type:     form.Type,
		AppId:    form.AppId,
		Sort:     form.Sort,
		Redirect: form.Redirect,
		Status:   form.Status,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

/*创建banner*/
func (ctl BannerController) Update(c *gin.Context) (gin.H, error) {
	var form UpdateBannerForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultBannerService.Update(srv_types.UpdateBannerDTO{
		Id:       form.Id,
		Name:     form.Name,
		ImageUrl: form.ImageUrl,
		Scene:    form.Scene,
		Type:     form.Type,
		AppId:    form.AppId,
		Sort:     form.Sort,
		Redirect: form.Redirect,
		Status:   form.Status,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

/*根据分页获取渠道列表*/
func (ctl BannerController) GetPageList(c *gin.Context) (gin.H, error) {
	var form GetBannerPageForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total, err := ctl.service.GetBannerPageList(srv_types.GetPageBannerDTO{
		Status: form.Status,
		Scene:  form.Scene,
		Name:   form.Name,
		Offset: form.Offset(),
		Limit:  form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}
