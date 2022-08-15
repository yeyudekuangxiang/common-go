package business

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/business/businesstypes"
	"mio/internal/pkg/core/context"
	sbusiness "mio/internal/pkg/service/business"
	"mio/internal/pkg/util/apiutil"
)

var DefaultAreaController = AreaController{}

type AreaController struct{}

func (AreaController) CityProvinceList(ctx *gin.Context) (gin.H, error) {
	form := businesstypes.CityProvinceForm{}

	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	areaSrv := sbusiness.NewAreaService(context.NewMioContext(context.WithContext(ctx)))
	list, err := areaSrv.GroupCityProvinceList(sbusiness.CityProvinceListDTO{
		Search: form.Search,
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list": list,
	}, nil
}
