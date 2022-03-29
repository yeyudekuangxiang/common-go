package service

import (
	"context"
	"errors"
	"fmt"
	"mio/config"
	"mio/core/app"
	"mio/internal/util"
	"mio/model/entity"
	"time"
)

var DefaultOCRService = NewOCRService()

func NewOCRService() OCRService {
	return OCRService{}
}

type OCRService struct {
}

// OCRForGm 素食打卡
func (u OCRService) OCRForGm(openid string, src string) error {
	res := util.OCRPush(src)
	var orderNo, fee string
	if res.WordsResultNum == 0 {
		return errors.New("无法识别此账单,请重新上传,谢谢")
	}
	for k, v := range res.WordsResult {
		if v.Words == "收款:" {
			fee = res.WordsResult[k+1].Words
		}
		if v.Words == "联系电话:021-62333696" {
			orderNo = res.WordsResult[k+1].Words
		}
	}
	fmt.Println("素食打卡账单:" + orderNo + ":" + fee)
	cmd := app.Redis.SetNX(context.Background(), config.RedisKey.Lock+orderNo, "a", 36500*time.Second)
	if !cmd.Val() {
		fmt.Println(config.RedisKey.Lock + orderNo + "重复扫描素食打卡")
		return errors.New("重复扫描素食打卡账单")
	}
	//发放积分
	_, err := DefaultPointTransactionService.Create(CreatePointTransactionParam{
		OpenId:       openid,
		Value:        100,
		Type:         entity.POINT_ADJUSTMENT,
		AdditionInfo: `{"素食打卡":"` + orderNo + `"}`,
	})

	return err
}
