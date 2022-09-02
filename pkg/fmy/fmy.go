package Fmy

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util"
)

type Fmy struct {
	ctx         *context.MioContext
	redis       *redis.Client
	appId       string
	appSecret   string
	platformKey string
}

func NewFmy(context *context.MioContext, client *redis.Client, appId, appSecret, platformKey string) *Fmy {
	return &Fmy{
		ctx:         context,
		redis:       client,
		appId:       appId,
		appSecret:   appSecret,
		platformKey: platformKey,
	}
}

func (f Fmy) GetSign(ch string, data []byte) (sign string, err error) {
	//查询 渠道信息
	scene := service.DefaultBdSceneService.FindByCh(ch)
	if scene.Key == "" || scene.Key == "e" {
		return "", errors.New("渠道查询失败")
	}
	rand1, rand2 := util.Rand4Number(), util.Rand4Number()
	fmt.Println(rand1, rand2)
	verifyData := rand1 + f.platformKey + string(data) + f.appSecret + rand2
	sign = rand1 + verifyData[6:27] + rand2
	fmt.Println(sign)
	return sign, nil
}
