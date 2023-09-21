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
	if topic.Status != 3 {
		return nil
	}
	ctx := context.Background()
	activeActivity, err := app.RpcService.ActivityRpcSrv.ActiveActivity(ctx, &activity.ActiveActivityReq{
		ActivityCode: "seeking_store",
	})
	if err != nil {
		return err
	}
	if !activeActivity.ActivityExist {
		return nil
	}
	if !activeActivity.ActivityRuleExist {
		return nil
	}

	startTime := activeActivity.ActiveActivity.ActivityRule.StartTime //毫秒
	endTime := activeActivity.ActiveActivity.ActivityRule.EndTime     //毫秒
	createTime := topic.CreatedAt.UnixMilli()
	if createTime < startTime || createTime > endTime {
		return errno.ErrTimeout.WithMessage("帖子不在当期活动时间内")
	}
	marshal, err := json.Marshal(topic)
	if err != nil {
		return err
	}
	err = app.QueueProduct.Publish(marshal, []string{routerkey.TopicSeekingStore}, rabbitmq.WithPublishOptionsExchange("lvmio"))
	if err != nil {
		return err
	}
	return nil
}
