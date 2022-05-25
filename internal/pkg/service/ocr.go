package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/httputil"
	"time"
)

var DefaultOCRService = NewOCRService()

func NewOCRService() OCRService {
	return OCRService{}
}

type OCRService struct {
}

// OCRForGm 素食打卡
func (srv OCRService) OCRForGm(openid string, src string) error {
	res := util.OCRPush(src)
	var orderNo, fee string

	for k, v := range res.WordsResult {
		if v.Words == "收款:" {
			fee = res.WordsResult[k+1].Words
		}
		if v.Words == "联系电话:021-62333696" {
			orderNo = res.WordsResult[k+1].Words
		}
	}
	if orderNo == "" || fee == "" {
		return errors.New("无法识别此账单,请重新上传,谢谢")
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
func (srv OCRService) Scan(imgUrl string) ([]string, error) {
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage?access_token=24.6157c4c9729181acc1bac04d6bd5ecbe.2592000.1650680140.282335-25833266"
	body, err := httputil.PostMapFrom(url, map[string]string{"url": "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/test/ocr/coffee_cup/greencat/65e2166e-708f-4298-921d-d0b314970a0f.jpg"})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body), err)
	var o OCRResult
	err = json.Unmarshal(body, &o)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for _, word := range o.WordsResult {
		results = append(results, word.Words)
	}
	return results, nil
}
