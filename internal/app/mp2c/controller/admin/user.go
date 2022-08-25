package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/wxamp"
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
			Partner:  entity.Partner(form.Partner),
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
		Partner:  form.Partner,
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
