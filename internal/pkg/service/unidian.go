package service

import (
	"context"
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util/unidian"
	"mio/pkg/errno"
	"time"
)

var DefaultUnidianService = NewUnidianService()

func NewUnidianService() UnidianService {
	return UnidianService{}
}

type UnidianService struct {
}

// SendPrize 发放奖励
func (u UnidianService) SendPrize(typeId string, mobile string, activityType string) error {
	//检测重复
	cmd := app.Redis.SetNX(context.Background(), activityType+mobile, "a", 3650*time.Second)
	if !cmd.Val() {
		fmt.Println(activityType + mobile + "重复充值")
		return errno.ErrCommon.WithMessage("正在充值,请稍等")
	}
	unidian.CouponOfUnidian(typeId, mobile, activityType+mobile)
	return nil
}