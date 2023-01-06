package cron

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/community"
)

func kumiaoCron() {
	//文章列表 每天一次更新热度
	AddFunc("@hourly", community.NewTopicService(context.NewMioContext()).ZAddTopic)
}
