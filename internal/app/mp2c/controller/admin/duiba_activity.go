package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultDuiBaActivityController = NewDuiBaActivityController(service.DefaultDuiBaActivityService)

func NewDuiBaActivityController(channel service.DuiBaActivityService) DuiBaActivityController {
	return DuiBaActivityController{service: channel}
}

type DuiBaActivityController struct {
	service service.DuiBaActivityService
}

//Create 创建
func (ctl DuiBaActivityController) Create(c *gin.Context) (gin.H, error) {
	var form CreateDuiBaActivityForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultDuiBaActivityService.Create(srv_types.CreateDuiBaActivityDTO{
		Name:        form.Name,
		ActivityUrl: form.ActivityUrl,
		Cid:         form.Cid,
		Type:        form.Type,
		ActivityId:  form.ActivityId,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

/*创建banner*/
func (ctl DuiBaActivityController) Update(c *gin.Context) (gin.H, error) {
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

/*
根据分页获取渠道列表
func (ctl DuiBaActivityController) GetPageList(c *gin.Context) (gin.H, error) {
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
*/
