package cron

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/kumiaoCommunity"
)

func kumiaoCron() {
	//文章列表 每天一次更新热度
	AddFunc("@hourly", kumiaoCommunity.NewTopicService(context.NewMioContext()).ZAddTopic)
}
