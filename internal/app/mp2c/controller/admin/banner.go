package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"strconv"
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
	insertDate := srv_types.CreateBannerDTO{
		Scene:    form.Scene,
		Type:     form.Type,
		Status:   form.Status,
		Name:     form.Name,
		ImageUrl: form.ImageUrl,
		AppId:    form.AppId,
		Sort:     form.Sort,
		Redirect: form.Redirect,
		Display:  form.Display,
	}
	if form.Scene == "event" {
		err := ctl.CreateSaas(insertDate)
		if err != nil {
			return nil, err
		}
		return gin.H{}, nil
	}
	err := service.DefaultBannerService.Create(insertDate)
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

type ListSaasResponse struct {
	Code       int         `json:"code"`
	DetailCode interface{} `json:"detailCode"`
	Message    string      `json:"message"`
	Data       struct {
		List []struct {
			Id       string `json:"id"`
			ImageUrl string `json:"imageUrl"`
			Name     string `json:"name"`
			Redirect string `json:"redirect"`
			Scene    string `json:"scene"`
			Sort     int    `json:"sort"`
			Status   int    `json:"status"`
			Type     string `json:"type"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
	Page interface{} `json:"page"`
}

type CreateSaasResponse struct {
	Code       int         `json:"code"`
	DetailCode interface{} `json:"detailCode"`
	Message    string      `json:"message"`
	Data       string      `json:"data"`
	Page       interface{} `json:"page"`
}

func (ctl BannerController) CreateSaas(insertDate srv_types.CreateBannerDTO) error {
	body, err := httptool.PostJson(config.Config.Saas.Domain+"/api/mp2c/spu/public-welfare/banner/edit", insertDate)
	if err != nil {
		app.Logger.Errorf("banner同步Saas系统出错: post error %s", err.Error())
		return err
	}
	response := CreateSaasResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return err
	}
	if response.Code != 200 {
		return errno.ErrCommon.WithErrMessage(response.Message)
	}
	return nil
}

func (ctl BannerController) UpdateSaas(insertDate srv_types.UpdateBannerDTO) error {
	body, err := httptool.PostJson(config.Config.Saas.Domain+"/api/mp2c/spu/public-welfare/banner/edit", insertDate)
	if err != nil {
		app.Logger.Errorf("banner同步Saas系统出错: post error %s", err.Error())
		return err
	}
	response := CreateSaasResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return err
	}
	if response.Code != 200 {
		return errno.ErrCommon.WithErrMessage(response.Message)
	}
	return nil
}

func (ctl BannerController) GetPageListSaas(insertDate srv_types.GetPageBannerDTO) ([]entity.Banner, int64, error) {
	list := make([]entity.Banner, 0)
	url := fmt.Sprintf("%s/api/mp2c/spu/public-welfare/banner/list?pageSize=%d&page=%d&name=%s&status=%d", config.Config.Saas.Domain, insertDate.Limit, insertDate.Offset, insertDate.Name, insertDate.Status)
	body, err := httptool.Get(url)
	if err != nil {
		app.Logger.Errorf("banner同步Saas系统出错: post error %s", err.Error())
		return list, 0, err
	}
	response := ListSaasResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Unmarshal body: %s\n", err.Error())
		return list, 0, err
	}
	if response.Code != 200 {
		return list, 0, errno.ErrCommon.WithErrMessage(response.Message)
	}
	if response.Data.Total == 0 {
		return list, 0, nil
	}
	for _, s := range response.Data.List {
		id, _ := strconv.ParseInt(s.Id, 10, 64)
		list = append(list, entity.Banner{
			ID:       id,
			Name:     s.Name,
			ImageUrl: s.ImageUrl,
			Scene:    entity.BannerScene(s.Scene),
			Type:     entity.BannerType(s.Type),
			Sort:     s.Sort,
			Redirect: s.Redirect,
			Status:   entity.BannerStatus(s.Status),
		})
	}
	return list, int64(response.Data.Total), nil

}

/*创建banner*/
func (ctl BannerController) Update(c *gin.Context) (gin.H, error) {
	var form UpdateBannerForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	updateDate := srv_types.UpdateBannerDTO{
		Id:       form.Id,
		Name:     form.Name,
		ImageUrl: form.ImageUrl,
		Scene:    form.Scene,
		Type:     form.Type,
		AppId:    form.AppId,
		Sort:     form.Sort,
		Redirect: form.Redirect,
		Status:   form.Status,
	}
	if form.Scene == "event" {
		err := ctl.UpdateSaas(updateDate)
		if err != nil {
			return nil, err
		}
		return gin.H{}, nil
	}
	err := service.DefaultBannerService.Update(updateDate)
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
	ParamsDate := srv_types.GetPageBannerDTO{
		Name:    form.Name,
		Scene:   form.Scene,
		Status:  form.Status,
		Offset:  form.Offset(),
		Limit:   form.Limit(),
		Display: form.Display,
	}
	if form.Scene == "event" {
		list, total, err := ctl.GetPageListSaas(ParamsDate)
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
	list, total, err := ctl.service.GetBannerPageList(ParamsDate)
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
