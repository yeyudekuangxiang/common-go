package service

import (
	"fmt"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	util2 "mio/internal/pkg/util/encrypt"
)

var DefaultBdSceneService = BdSceneService{}

type BdSceneService struct {
}

func (srv BdSceneService) FindByCh(ch string) *entity.BdScene {
	item := repository.DefaultBdSceneRepository.FindByCh(ch)
	return &item
}

func (srv BdSceneService) SceneToType(ch string) entity.PointTransactionType {
	return repository.DefaultBdSceneRepository.SceneToType(ch)
}

func (srv BdSceneService) SceneToCarbonType(ch string) entity.CarbonTransactionType {
	return repository.DefaultBdSceneRepository.SceneToCarbonType(ch)
}

func (srv BdSceneService) CheckSign(mobile string, outTradeNo string, total string, sign string, scene *entity.BdScene) bool {
	str := scene.Ch + "#" + mobile + "#" + outTradeNo + "#" + total + "#" + scene.Key
	fmt.Println("localSignStr", str)
	localSign := util2.Md5(str)
	fmt.Println("CheckSign", localSign, sign)
	return localSign == sign
}
