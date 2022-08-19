package util

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"mio/internal/pkg/core/app"
	"sync"
	"testing"
)

func TestSnowflakeNodeID(t *testing.T) {
	app.Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := app.Redis.Ping(context.Background()).Err(); err != nil {
		log.Println("redis未链接 跳过雪花id测试")
		return
	}
	wait := sync.WaitGroup{}
	data := make(map[int64]int64)
	m := sync.Mutex{}
	for i := 0; i < 512; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for j := 0; j < 20; j++ {
				id, err := app.Redis.Eval(context.Background(), snowflakeLua, []string{"snowflake"}).Int64()
				assert.Equal(t, nil, err, "生成雪花id异常")
				if err != nil {
					assert.LessOrEqual(t, id, int64(1024), "生成雪花id范围异常")
				}
				m.Lock()
				data[id] = data[id] + 1
				assert.LessOrEqual(t, data[id], int64(100), "生成雪花id重复")
				m.Unlock()
			}
		}()
	}

	wait.Wait()

	fmt.Printf("id列表%+v\n", data)

}
