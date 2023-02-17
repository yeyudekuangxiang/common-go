package cron

import (
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/service/community"
)

func kumiaoCron() {
	//文章列表 每小时更新热度
	AddFunc("@hourly", community.NewTopicService(context.NewMioContext()).ZAddTopic)
	//文章列表 每周日零点更新
	AddFunc("@weekly", community.NewTopicService(context.NewMioContext()).SetWeekTopic)
}
