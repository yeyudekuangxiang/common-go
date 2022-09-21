package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultOCRController = OCRController{}

type OCRController struct {
}

func (OCRController) GmTicket(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	form := GetOCRForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultOCRService().OCRForGm(user.OpenId, user.Risk, form.Src)
	return gin.H{}, err
}
