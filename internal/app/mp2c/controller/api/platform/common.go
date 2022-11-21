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

//天津地铁获取用户状态

func (ctr CommonController) GetTjMetroTicketStatus(c *gin.Context) (gin.H, error) {
	//user := apiutil.GetAuthUser(c)
	user, _, _ := service.DefaultUserService.GetUser(repository.GetUserBy{OpenId: "oMD8d5CPOCCTAzfohzl_3t7ZBBB0"})
	serviceNew := tianjin_metro.NewTianjinMetroService(context.NewMioContext())
	err := serviceNew.GetTjMetroTicketStatus(user.OpenId, user.ID, user.ChannelId, 1000)
	if err != nil {
		return gin.H{}, err
	}
	return gin.H{}, nil
}
