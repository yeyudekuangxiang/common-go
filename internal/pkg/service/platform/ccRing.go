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
func (c ccRing) CallBack(userInfo *entity.User, degree float64) {
	sceneUser := repository.DefaultBdSceneUserRepository.FindPlatformUserByOpenId(userInfo.OpenId)
	if sceneUser.ID != 0 && sceneUser.PlatformKey == "ccring" {
		scene := repository.DefaultBdSceneRepository.FindByCh("ccring")
		if scene.ID == 0 {
			app.Logger.Errorf("回调光环错误:%s", "未设置scene")
			return
		}
		url := scene.Domain + "/api/cc-ring/external/ev-charge"
		authToken := httputil.HttpWithHeader("Authorization", "dsaflsdkfjxcmvoxiu123moicuvhoi123")
		queryParams := ccRingReqParams{
			MemberId:       sceneUser.PlatformUserId,
			DegreeOfCharge: degree,
		}
		_, err := httputil.PostJson(url, queryParams, authToken)
		if err != nil {
			app.Logger.Errorf("回调光环错误:%s", err.Error())
			return
		}
		return
	}
}
