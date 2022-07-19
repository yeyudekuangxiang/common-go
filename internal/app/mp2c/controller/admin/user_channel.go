package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultUserChannelController = NewUserChannelController(service.DefaultUserChannelService)

func NewUserChannelController(channel service.UserChannelService) UserChannelController {
	return UserChannelController{service: channel}
}

type UserChannelController struct {
	service service.UserChannelService
}

/*创建渠道*/
func (ctl UserChannelController) Create(c *gin.Context) (gin.H, error) {
	var form CreateUserChannelForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := ctl.service.Create(&entity.UserChannel{
		Cid:     form.Cid,
		Pid:     form.Pid,
		Name:    form.Name,
		Code:    form.Code,
		Company: form.Company,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

/*根据渠道id，更新渠道信息*/
func (ctl UserChannelController) UpdateByCid(c *gin.Context) (gin.H, error) {
	var form CreateUserChannelForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := ctl.service.UpdateUserChannel(&entity.UserChannel{
		Cid:     form.Cid,
		Pid:     form.Pid,
		Name:    form.Name,
		Code:    form.Code,
		Company: form.Company,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

/*根据分页获取渠道列表*/
func (ctl UserChannelController) GetPageList(c *gin.Context) (gin.H, error) {
	var form GetUserChannelPageForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total, err := ctl.service.GetUserChannelPageList(repository.GetUserChannelPageListBy{
		Cid:    form.Cid,
		Pid:    form.Pid,
		Code:   form.Code,
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
