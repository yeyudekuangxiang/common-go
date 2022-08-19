package service

import (
	"errors"
	"fmt"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	util2 "mio/internal/pkg/util/encrypt"
	"strings"
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

func (srv BdSceneService) CheckSign(mobile string, outTradeNo string, total string, sign string, scene *entity.BdScene) bool {
	str := scene.Ch + "#" + mobile + "#" + outTradeNo + "#" + total + "#" + scene.Key
	fmt.Println("localSignStr", str)
	localSign := util2.Md5(str)
	fmt.Println("CheckSign", localSign, sign)
	return localSign == sign
}

func (srv BdSceneService) CheckWhiteList(ip, ch string) error {
	result := repository.DefaultBdSceneRepository.FindByCh(ch)
	if result.WhiteIp != "e" {
		items := strings.Split(result.WhiteIp, ",")
		for _, item := range items {
			if ip == item {
				return nil
			}
		}
		return errors.New("非白名单ip")
	}
	return nil
}
