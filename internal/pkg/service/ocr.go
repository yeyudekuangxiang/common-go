package service

import (
	"context"
	"errors"
	"fmt"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/baidu"
	"time"
)

var DefaultOCRService OCRService

func InitDefaultOCRService() {
	DefaultOCRService = OCRService{
		imageClient: &baidu.ImageClient{
			AccessToken: baidu.NewAccessToken(baidu.AccessTokenConfig{
				RedisClient: app.Redis,
				Prefix:      config.RedisKey.BaiDu,
				AppKey:      config.Config.BaiDu.AppKey,
				AppSecret:   config.Config.BaiDu.AppSecret,
			}),
		},
	}
}

type OCRService struct {
	imageClient *baidu.ImageClient
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

	pointTranService := NewPointService(mioctx.NewMioContext())
	//发放积分
	_, err := pointTranService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openid,
		ChangePoint:  100,
		BizId:        util.UUID(),
		Type:         entity.POINT_ADJUSTMENT,
		AdditionInfo: `{"素食打卡":"` + orderNo + `"}`,
	})

	return err
}
func (srv OCRService) Scan(imgUrl string) ([]string, error) {
	rest, err := srv.imageClient.WebImage(baidu.WebImageParam{
		ImageUrl: imgUrl,
	})
	fmt.Printf("%+v %+v\n", rest, err)
	if err != nil {
		return nil, err
	}
	if !rest.IsSuccess() {
		return nil, errors.New(rest.ErrorDescription)
	}

	results := make([]string, 0)
	for _, word := range rest.WordsResult {
		results = append(results, word.Words)
	}
	return results, nil
}
