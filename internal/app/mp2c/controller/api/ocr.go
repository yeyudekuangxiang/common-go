package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

var DefaultOCRController = OCRController{}

type OCRController struct {
}

func (OCRController) GmTicket(c *gin.Context) (gin.H, error) {
	user := util.GetAuthUser(c)
	form := GetOCRForm{}
	if err := util.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultOCRService.OCRForGm(user.OpenId, form.Src)
	return gin.H{}, err
}
