package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/wxamp"
	"mio/pkg/errno"
	"strings"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (ctr UserController) List(c *gin.Context) (gin.H, error) {
	form := UserPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total := service.DefaultUserService.GetUserPageListBy(repository.GetUserPageListBy{
		Limit:  form.Limit(),
		Offset: form.Offset(),
		User: repository.GetUserListBy{
			Mobile:   form.Mobile,
			UserId:   form.ID,
			Status:   form.State,
			Nickname: form.Nickname,
			Position: entity.UserPosition(form.Position),
			Partners: entity.Partner(form.Partners),
		},
		OrderBy: "id desc",
	})
	return gin.H{
		"users":    list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (ctr UserController) Detail(c *gin.Context) (gin.H, error) {
	var form IDForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user, err := service.DefaultUserService.GetUserById(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}

func (ctr UserController) Update(c *gin.Context) (gin.H, error) {
	var form UpdateUser
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	if err := service.DefaultUserService.UpdateUserInfo(service.UpdateUserInfoParam{
		UserId:   form.ID,
		Status:   form.Status,
		Position: form.Position,
		Partners: form.Partners,
		Auth:     form.Auth,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr UserController) PositionList(c *gin.Context) (gin.H, error) {
	var position []string
	position = append(position, "yellow", "blue", "ordinary")
	return gin.H{
		"position": position,
	}, nil
}

func (ctr UserController) UpdateUserRisk(c *gin.Context) (gin.H, error) {
	i := 0
	for {
		list, _ := service.DefaultUserService.GetUserPageListBy(repository.GetUserPageListBy{
			Limit:   10,
			Offset:  i,
			OrderBy: "id desc",
			User:    repository.GetUserListBy{Risk: -1},
		})
		if len(list) == 0 {
			break
		}
		i += len(list)
		var ids []string
		for _, v := range list {
			if strings.Contains(v.OpenId, "oy_") {
				ids = append(ids, v.OpenId)
			}
		}

		//openid 一次最多传十个
		cas := wxamp.BatchGetUserRiskCase(ids)
		if cas == nil {
			continue
		}
		i -= len(cas.List)
		//保存risk
		for _, v := range list {
			for _, c := range cas.List {
				if v.OpenId == c.Openid {
					fmt.Println("checkopenid", v.ID, c.RiskRank, v.OpenId)
					service.DefaultUserService.UpdateUserRisk(service.UpdateUserRiskParam{
						UserId: v.ID,
						Risk:   c.RiskRank,
					})

				}
			}
		}
		fmt.Println("risk i ", i)

	}

	return nil, nil
}

//用户列表 test

func (ctr UserController) ListRisk(c *gin.Context) (gin.H, error) {
	form := UserPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, total := service.NewUserRiskService(context.NewMioContext()).GetUserRiskPageListBy(repository.GetUserPageListBy{
		Limit:  form.Limit(),
		Offset: form.Offset(),
		User: repository.GetUserListBy{
			Mobile:   form.Mobile,
			UserId:   form.ID,
			Status:   form.State,
			Nickname: form.Nickname,
		},
		OrderBy: "id desc",
	})
	return gin.H{
		"users":    list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

//risk统计分类

func (ctr UserController) RiskStatistics(c *gin.Context) (gin.H, error) {
	list := service.NewUserRiskService(context.NewMioContext()).GetUserRiskStatisticst()
	return gin.H{
		"date": list,
	}, nil
}

//更新用户风险等级

func (ctr UserController) UpdateRisk(c *gin.Context) (gin.H, error) {
	var form UpdateUserRisk
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	idsSlice := strings.Split(form.Ids, ",")
	var UserIdSlice, PhoneSlice, OpenIdSlice []string

	if len(idsSlice) == 0 {
		return nil, errno.ErrCommon.WithMessage("请输出要提交的id")
	}
	switch form.Type {
	case 1:
		UserIdSlice = idsSlice
		break
	case 2:
		PhoneSlice = idsSlice
		break
	case 3:
		OpenIdSlice = idsSlice
		break
	default:
		break
	}

	if err := service.NewUserRiskService(context.NewMioContext()).BatchUpdateUserRisk(service.UpdateRiskParam{
		UserIdSlice: UserIdSlice,
		PhoneSlice:  PhoneSlice,
		OpenIdSlice: OpenIdSlice,
		Risk:        form.Risk,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
