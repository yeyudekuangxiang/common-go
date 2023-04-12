package community

import (
	"context"
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.miotech.com/miotech-application/backend/mp2c-micro/app/activity/cmd/rpc/activity/activity"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/types/message/communitymsg"
	"mio/internal/pkg/queue/types/routerkey"
	"mio/pkg/errno"
)

//SeekingStore 探店
func SeekingStore(topic communitymsg.Topic) error {
	//查询是否为探店标签
	var flag bool
	for _, tag := range topic.Tags {
		if tag.Name == "城市碳秘" {
			flag = true
		}
	}
	if !flag {
		return nil
	}
	ctx := context.Background()
	rule, err := app.RpcService.ActivityRpcSrv.ActiveRule(ctx, &activity.ActiveRuleReq{
		Code: "seeking_store",
	})
	if err != nil {
		return err
	}
	if !rule.GetExist() {
		return errno.ErrRecordNotFound
	}

	marshal, err := json.Marshal(topic)
	if err != nil {
		app.Logger.Errorf("[城市碳秘] json_encode 错误: %+v\n", err.Error())
		return err
	}
	err = app.QueueProduct.Publish(marshal, []string{routerkey.TopicSeekingStore}, rabbitmq.WithPublishOptionsExchange("lvmio"))
	if err != nil {
		app.Logger.Errorf("[城市碳秘] 错误:%+v\n", err.Error())
		return err
	}
	return nil
}
