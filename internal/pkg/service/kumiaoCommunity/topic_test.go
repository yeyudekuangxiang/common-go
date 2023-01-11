package kumiaoCommunity

import (
	"fmt"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/repository"
	"strings"
	"testing"
)

//func Test_NLOC(t *testing.T) {
//	initialize.Initialize("/Users/yunfeng/Documents/workspace/mp2c-go/config.ini")
//	list := make([]entity.Topic, 0)
//	app.DB.Model(&entity.Topic{}).Order("updated_at desc").Limit(20).Find(&list)
//	timer := time.NewTimer(time.Second * 10)
//	quit := make(chan struct{})
//
//	defer timer.Stop()
//	go func() {
//		<-timer.C
//		close(quit)
//	}()
//
//	go func() {
//		n := util.NewHot()
//		for _, item := range list {
//			score := n.Hot(int64(item.SeeCount), item.LikeCount, item.CommentCount, item.CreatedAt.Time)
//			fmt.Printf("score: %f\n", score)
//		}
//	}()
//
//	for {
//		<-quit
//		return
//	}
//}

func TestSetHotTopic(t *testing.T) {
	initialize.Initialize("/Users/yunfeng/Documents/workspace/mp2c-go/config.ini")
	//topicSercice := NewTopicService(context.NewMioContext()).SetWeekTopic()
	topics, err := repository.NewTopicModel(context.NewMioContext()).GetImportTopic()
	if err != nil {
		fmt.Printf("err1 is %s", err.Error())
	}

	if len(topics) == 0 {
		fmt.Printf("err2 is %s", "not found topic")
	}

	uMapToTopic := make(map[int64]struct{}, 0) // 每个用户一篇文章
	uSliceToTopic := make([]int64, 0)
	var j int64
	for _, item := range topics {
		if _, ok := uMapToTopic[item.UserId]; ok {
			continue
		}
		uMapToTopic[item.UserId] = struct{}{}
		uSliceToTopic = append(uSliceToTopic, item.Id)
		if j == 49 {
			break
		}
		j++
	}

	fmt.Printf("uSliceToTopic: %v", uSliceToTopic)
}

func TestSplit(t *testing.T) {
	keys := strings.Split("wechat", "_")
	fmt.Println(keys)
}
