package common

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/common"
	"mio/internal/pkg/util/apiutil"
)

var DefaultCityController = CityController{}

type CityController struct {
}

func (ctr CityController) List(c *gin.Context) (gin.H, error) {
	form := CityListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user := apiutil.GetAuthUser(c)
	citySrv := common.NewCityService(context.NewMioContext(context.WithContext(c.Request.Context())))

	list, err := citySrv.GetCityList(common.GetCityListParams{
		//CityCode:    form.CityCode,
		CityPidCode: form.CityPidCode,
	})

	if err != nil {
		return nil, err
	}

	return gin.H{
		"list": list,
	}, nil
}
