package util

import (
	"context"
	"github.com/bwmarrin/snowflake"
	uuid "github.com/iris-contrib/go.uuid"
	"mio/internal/pkg/core/app"
	"sync"
)

// UUID 生成uuid出现错误时会panic
func UUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

var sfInc *snowflake.Node
var onceSfInc sync.Once

const snowflakeLua = `local num = redis.pcall("INCR",KEYS[1]) 
if num > 1024 then 
redis.pcall("SET",KEYS[1],num-1024) 
return num-1024 
else 
return num 
end`

// SnowflakeID 生成雪花id
func SnowflakeID() (snowflake.ID, error) {
	var e error
	onceSfInc.Do(func() {
		nodeId, err := app.Redis.Eval(context.Background(), snowflakeLua, []string{"snowflake"}).Int64()
		if err != nil {
			e = err
			return
		}

		node, err := snowflake.NewNode(nodeId)
		if err != nil {
			e = err
			return
		}

		sfInc = node
	})
	if e != nil {
		return 0, e
	}

	return sfInc.Generate(), nil
}
