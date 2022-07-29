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
func (BadgeController) GetBadgeList(ctx *gin.Context) (gin.H, error) {

	user := apiutil.GetAuthUser(ctx)
	list, err := service.DefaultBadgeService.GetBadgePageList(user.OpenId)

	return gin.H{
		"list": list,
	}, err
}
func (BadgeController) UpdateBadgeIsNew(ctx *gin.Context) (gin.H, error) {

	form := UpdateBadgeIsNewForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(ctx)
	err := service.DefaultBadgeService.UpdateBadgeIsNew(user.OpenId, form.BadgeId, true)

	return nil, err
}
func (BadgeController) UploadOldBadgeImage(ctx *gin.Context) (interface{}, error) {
	form := UploadOLdBadgeImageForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	info, err := service.DefaultBadgeService.GetUploadOldBadgeSetting(form.BadgeId)
	if err != nil {
		return nil, err
	}
	return info, nil
}
