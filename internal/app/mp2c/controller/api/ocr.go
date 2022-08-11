package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/pkg/errno"
	"time"
)

var DefaultOCRController = OCRController{}

type OCRController struct {
}

func (OCRController) GmTicket(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	if !util.DefaultLock.Lock("GmTicket:"+user.OpenId, time.Second*5) {
		return nil, errno.ErrLimit.WithCaller()
	}
	if user.Risk > 2 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}
	form := GetOCRForm{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultOCRService.OCRForGm(user.OpenId, form.Src)
	return gin.H{}, err
}
