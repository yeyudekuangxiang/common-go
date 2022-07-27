package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model/entity"
)

var DefaultConstantController = NewConstantController()

func NewConstantController() ConstantController {
	return ConstantController{}
}

type ConstantController struct {
}

func (ctl ConstantController) List(c *gin.Context) (gin.H, error) {
	BannerCollectValueMap := map[string]interface{}{
		"bannerStatus": entity.BannerStatusMap,
		"bannerScene":  entity.BannerSceneMap,
		"BannerType":   entity.BannerTypeMap,
	}
	return gin.H{
		"banner": BannerCollectValueMap,
	}, nil
}
