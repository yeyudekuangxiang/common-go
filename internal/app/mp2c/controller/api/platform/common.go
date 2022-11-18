package platform

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/tianjin_metro"
)

var DefaultCommonController = CommonController{}

type CommonController struct {
}

func (ctr CommonController) SendMetro(c *gin.Context) (gin.H, error) {
	//user := apiutil.GetAuthUser(c)
	userV2, _, _ := service.DefaultUserService.GetUser(repository.GetUserBy{OpenId: "oMD8d5CPOCCTAzfohzl_3t7ZBBB0"})

	serviceNew := tianjin_metro.NewTianjinMetroService(context.NewMioContext())
	str, err := serviceNew.SendCoupon(1, 1, *userV2)
	if err != nil {
		return gin.H{}, nil
	}
	return gin.H{
		"str": str,
	}, nil

}

func (ctr CommonController) GetTjMetroTicketStatus(c *gin.Context) (gin.H, error) {
	//user := apiutil.GetAuthUser(c)
	userV2, _, _ := service.DefaultUserService.GetUser(repository.GetUserBy{OpenId: "oMD8d5CPOCCTAzfohzl_3t7ZBBB0"})

	serviceNew := tianjin_metro.NewTianjinMetroService(context.NewMioContext())
	str, err := serviceNew.GetTjMetroTicketStatus(userV2.OpenId)
	if err != nil {
		return gin.H{}, err
	}
	return gin.H{
		"str": str,
	}, nil

}
