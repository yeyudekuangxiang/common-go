package badge

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultBadgeController = BadgeController{}

type BadgeController struct {
}

func (BadgeController) UpdateBadgeImage(ctx *gin.Context) (gin.H, error) {
	form := UpdateBadgeImageForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(ctx)
	return nil, service.DefaultBadgeService.UpdateCertImage(user.OpenId, form.UploadCode, form.ImageUrl)
}
