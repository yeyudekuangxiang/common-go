package kumiaoCommunity

import (
	"fmt"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/util"
	"testing"
	"time"
)

func Test_NLOC(t *testing.T) {
	initialize.Initialize("/Users/yunfeng/Documents/workspace/mp2c-go/config.ini")
	list := make([]entity.Topic, 0)
	app.DB.Model(&entity.Topic{}).Order("updated_at desc").Limit(20).Find(&list)
	timer := time.NewTimer(time.Second * 10)
	quit := make(chan struct{})

	defer timer.Stop()
	go func() {
		<-timer.C
		close(quit)
	}()

	go func() {
		n := util.NewHot()
		for _, item := range list {
			score := n.Hot(int64(item.SeeCount), item.LikeCount, item.CommentCount, item.CreatedAt.Time)
			fmt.Printf("score: %f\n", score)
		}
	}()

	for {
		<-quit
		return
	}
}
