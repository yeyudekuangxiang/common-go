package service

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/common/tool/encrypt"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	util2 "mio/internal/pkg/util/encrypt"
	"mio/pkg/errno"
	"sort"
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

func (srv BdSceneService) CheckPreSign(key, sign string, params map[string]string) bool {
	var slice []string
	for k := range params {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	var signStr string
	for _, v := range slice {
		signStr += v + "=" + params[v] + "&"
	}
	signStr = strings.TrimRight(signStr, "&")
	fmt.Println(encrypt.Md5(key + signStr))
	return encrypt.Md5(key+signStr) == sign
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
		return errno.ErrCommon.WithMessage("非白名单ip")
	}
	return nil
}
