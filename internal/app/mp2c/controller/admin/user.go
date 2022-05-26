package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/wxamp"
	"strings"
)

var DefaultUserController = UserController{}

type UserController struct {
}

func (UserController) GetUserInfo(c *gin.Context) (gin.H, error) {
	var form GetUserForm
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user, err := service.DefaultUserService.GetUserById(form.Id)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"user": user,
	}, nil
}

func GetUserPageListBy(c *gin.Context) (gin.H, error) {
	form := UserPageListForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	list, count := service.DefaultUserService.GetUserPageListBy(repository.GetUserPageListBy{
		Limit:   10,
		Offset:  0,
		User:    repository.GetUserListBy{Mobile: form.Mobile},
		OrderBy: "id desc",
	})
	return gin.H{
		"users": list,
		"page":  count,
	}, nil
}

func UpdateUserRisk(c *gin.Context) (gin.H, error) {
	i := 0
	for {
		list, count := service.DefaultUserService.GetUserPageListBy(repository.GetUserPageListBy{
			Limit:   10,
			Offset:  i * 10,
			OrderBy: "id desc",
		})

		var ids []string
		for _, v := range list {
			if strings.Contains(v.OpenId, "oy_") {
				ids = append(ids, v.OpenId)
			}
		}

		//openid 一次最多传十个
		cas := wxamp.BatchGetUserRiskCase(ids)
		if cas == nil {
			i++
			continue
		}
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

		i++
		if i*10 >= int(count) {
			break
		}

	}

	return nil, nil
}
