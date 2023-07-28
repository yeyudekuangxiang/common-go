package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"net/url"
	"strings"
)

var DefaultDuiBaActivityController = NewDuiBaActivityController()

func NewDuiBaActivityController() DuiBaActivityController {
	return DuiBaActivityController{}
}

type DuiBaActivityController struct {
}

// Save 创建兑吧链接
func (ctl DuiBaActivityController) Create(c *gin.Context) (gin.H, error) {
	var form CreateDuiBaActivityForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	newUrl := strings.Replace(form.ActivityUrl, "https://go-api.miotech.com?dbredirect=", "", 1)
	url, errUrl := url.QueryUnescape(newUrl)
	if errUrl != nil {
		return nil, errno.ErrCommon.WithMessage("链接格式错误")
	}
	if find := strings.Contains(url, "http"); !find {
		url = "https:" + url
	}
	err := service.NewDuiBaActivityService(context.NewMioContext()).Create(srv_types.CreateDuiBaActivityDTO{
		Name:           form.Name,
		Cid:            form.Cid,
		Type:           form.Type,
		IsShare:        form.IsShare,
		IsPhone:        form.IsPhone,
		ActivityUrl:    url,
		ActivityId:     form.ActivityId,
		RiskLimit:      form.RiskLimit,
		BlackWhiteType: form.BlackWhiteType,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

// Update 修改兑吧链接
func (ctl DuiBaActivityController) Update(c *gin.Context) (gin.H, error) {
	var form UpdateDuiBaActivityForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	newUrl := strings.Replace(form.ActivityUrl, "https://go-api.miotech.com?dbredirect=", "", 1)
	url, errUrl := url.QueryUnescape(newUrl)
	if errUrl != nil {
		return nil, errno.ErrCommon.WithMessage("链接格式错误")
	}
	if find := strings.Contains(url, "http"); !find {
		url = "https:" + url
	}
	err := service.NewDuiBaActivityService(context.NewMioContext()).Update(srv_types.UpdateDuiBaActivityDTO{
		Name:        form.Name,
		ActivityUrl: url,
		Cid:         form.Cid,
		Type:        form.Type,
		ActivityId:  form.ActivityId,
		IsShare:     form.IsShare,
		Id:          form.Id,
		IsPhone:     form.IsPhone,
		RiskLimit:   form.RiskLimit,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

// GetPageList 根据分页获取渠道列表
func (ctl DuiBaActivityController) GetPageList(c *gin.Context) (gin.H, error) {
	var form GetDuiBaActivityPageForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total, err := service.NewDuiBaActivityService(context.NewMioContext()).GetPageList(srv_types.GetPageDuiBaActivityDTO{
		ActivityId: form.ActivityId,
		Type:       form.Type,
		Status:     entity.DuiBaActivityStatusYes,
		Cid:        form.Cid,
		Name:       form.Name,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	voList := make([]api_types.DuiBaActivityVO, 0)
	for _, activity := range list {
		/*res, err1, err2 := app.Weapp.GetQRCode(&weapp.QRCode{
			Path: fmt.Sprintf("/pages/duiba_v2/duiba/index?activityId=%s", activity.ActivityId),
		})
		if err2 != nil {
			return nil,err2
		}
		if err1 != nil {
			return nil,err1
		}
		*/
		urlVo := url.QueryEscape(activity.ActivityUrl)
		urlVo = "https://go-api.miotech.com?dbredirect=" + urlVo
		voList = append(voList, api_types.DuiBaActivityVO{
			ID:             activity.ID,
			Name:           activity.Name,
			Type:           activity.Type,
			Cid:            activity.Cid,
			IsShare:        activity.IsShare,
			IsPhone:        activity.IsPhone,
			ActivityId:     activity.ActivityId,
			ActivityUrl:    urlVo,
			RiskLimit:      activity.RiskLimit,
			CreatedAt:      activity.CreatedAt.Format("2006.01.02 15:04:05"),
			UpdatedAt:      activity.UpdatedAt.Format("2006.01.02 15:04:05"),
			BlackWhiteType: activity.BlackWhiteType,
		})
	}
	return gin.H{
		"list":     voList,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

// Delete 修改兑吧链接
func (ctl DuiBaActivityController) Delete(c *gin.Context) (gin.H, error) {
	var form DeleteDuiBaActivityForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	err := service.NewDuiBaActivityService(context.NewMioContext()).Delete(srv_types.DeleteDuiBaActivityDTO{
		Id: form.Id,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

// Show 根据分页获取渠道列表
func (ctl DuiBaActivityController) Show(c *gin.Context) (gin.H, error) {
	var form ShowDuiBaActivityForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	info, err := service.NewDuiBaActivityService(context.NewMioContext()).Show(srv_types.ShowDuiBaActivityDTO{
		Id: form.Id,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"info": info,
	}, nil
}
