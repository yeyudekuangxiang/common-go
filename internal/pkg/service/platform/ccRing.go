package platform

import (
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/util/httputil"
)

var DefaultCCRingService = NewCCRingService()

func NewCCRingService() *ccRing {
	return &ccRing{}
}

type ccRing struct {
}

//回调ccring
func (c ccRing) CallBack(userInfo *entity.User, degree float64, platformKey string, scene *entity.BdScene) {
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUser(userInfo.OpenId, platformKey)
	if sceneUser.ID != 0 {
		url := scene.Domain + "/api/cc-ring/external/ev-charge"
		authToken := httputil.HttpWithHeader("Authorization", "dsaflsdkfjxcmvoxiu123moicuvhoi123")
		queryParams := CcRingReqParams{
			MemberId:       sceneUser.PlatformUserId,
			DegreeOfCharge: degree,
		}
		_, err := httputil.PostJson(url, queryParams, authToken)
		if err != nil {
			app.Logger.Errorf("回调光环错误:%s", err.Error())
		}
	}
}
